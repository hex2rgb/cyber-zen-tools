package commands

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// newCompressCommand 创建图片压缩命令
func newCompressCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compress",
		Short: "压缩图片文件",
		Long: `压缩图片文件，支持多种格式。

压缩策略：
  1. 优先保证图片质量（无损压缩）
  2. 按指定比率缩小图片尺寸
  3. 自动优化文件大小

支持的格式：
  - JPEG (.jpg, .jpeg): 质量优化 + 尺寸调整
  - PNG (.png): 无损压缩 + 尺寸调整  
  - GIF (.gif): 尺寸调整
  - 其他格式: 直接复制

特性：
  - 自动添加时间戳避免覆盖
  - 保持原文件扩展名
  - 支持相对路径和绝对路径
  - 自动创建目标目录

用法:
  cyber-zen compress --src "源文件或文件夹" --dist "目标路径" --rate "压缩比率"

参数:
  --src   源文件或文件夹路径（必需）
  --dist  目标路径（可选，默认当前目录）
  --rate  压缩比率 0.1-1.0（可选，默认0.8）

示例:
  cyber-zen compress --src "images/" --dist "compressed/" --rate 0.7
  cyber-zen compress --src "photo.jpg" --rate 0.5`,
		RunE: func(cmd *cobra.Command, args []string) error {
			src, _ := cmd.Flags().GetString("src")
			dist, _ := cmd.Flags().GetString("dist")
			rate, _ := cmd.Flags().GetFloat64("rate")
			
			return runCompress(src, dist, rate)
		},
	}

	// 添加参数
	cmd.Flags().String("src", "", "源文件或文件夹路径")
	cmd.Flags().String("dist", "", "目标路径（可选）")
	cmd.Flags().Float64("rate", 0.8, "压缩比率 0.1-1.0（可选）")
	
	// 标记必需参数
	cmd.MarkFlagRequired("src")

	return cmd
}

// runCompress 执行图片压缩
func runCompress(src, dist string, rate float64) error {
	color.Green("开始压缩图片...")
	color.Cyan("源路径: %s", src)
	color.Cyan("目标路径: %s", dist)
	color.Cyan("压缩比率: %.2f", rate)

	// 验证压缩比率
	if rate < 0.1 || rate > 1.0 {
		return fmt.Errorf("压缩比率必须在 0.1 到 1.0 之间")
	}

	// 获取绝对路径
	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return fmt.Errorf("获取源路径失败: %v", err)
	}

	// 检查源路径是否存在
	if _, err := os.Stat(srcAbs); os.IsNotExist(err) {
		return fmt.Errorf("源路径不存在: %s", srcAbs)
	}

	// 处理目标路径
	if dist == "" {
		// 如果目标路径为空，在当前目录创建带时间戳的文件夹
		timestamp := time.Now().Format("20060102_150405")
		dist = fmt.Sprintf("compressed_%s", timestamp)
	}
	distAbs, err := filepath.Abs(dist)
	if err != nil {
		return fmt.Errorf("获取目标路径失败: %v", err)
	}

	// 检查目标路径是否已经包含时间戳，如果没有则添加
	distWithTimestamp := distAbs
	if !containsTimestamp(distAbs) {
		timestamp := time.Now().Format("20060102_150405")
		distWithTimestamp = addTimestampToPath(distAbs, timestamp)
	}

	// 创建目标目录
	if err := os.MkdirAll(filepath.Dir(distWithTimestamp), 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 判断是文件还是目录
	fileInfo, err := os.Stat(srcAbs)
	if err != nil {
		return fmt.Errorf("获取源文件信息失败: %v", err)
	}

	if fileInfo.IsDir() {
		// 压缩目录
		return compressDirectory(srcAbs, distWithTimestamp, rate)
	} else {
		// 压缩单个文件：确保目标文件名包含原文件扩展名
		originalExt := filepath.Ext(srcAbs)
		originalName := filepath.Base(srcAbs)
		nameWithoutExt := strings.TrimSuffix(originalName, originalExt)
		
		// 如果目标路径是目录，在目录中创建带时间戳的文件
		if distAbs == distWithTimestamp {
			// 目标路径没有扩展名，说明是目录，需要添加文件名
			distWithTimestamp = filepath.Join(distWithTimestamp, fmt.Sprintf("%s_%s%s", nameWithoutExt, time.Now().Format("20060102_150405"), originalExt))
		}
		
		// 确保目标目录存在
		if err := os.MkdirAll(filepath.Dir(distWithTimestamp), 0755); err != nil {
			return fmt.Errorf("创建目标目录失败: %v", err)
		}
		
		return compressFile(srcAbs, distWithTimestamp, rate)
	}
}

// addTimestampToPath 在路径中添加时间戳
func addTimestampToPath(path, timestamp string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	
	if ext != "" {
		// 文件：保持原扩展名
		return filepath.Join(dir, fmt.Sprintf("%s_%s%s", name, timestamp, ext))
	} else {
		// 目录：添加时间戳
		return filepath.Join(dir, fmt.Sprintf("%s_%s", name, timestamp))
	}
}

// containsTimestamp 检查路径是否已经包含时间戳
func containsTimestamp(path string) bool {
	pattern := regexp.MustCompile(`_\d{8}_\d{6}$`) // 匹配 _YYYYMMDD_HHMMSS 格式
	return pattern.MatchString(filepath.Base(path))
}

// compressDirectory 压缩目录中的所有图片
func compressDirectory(srcDir, distDir string, rate float64) error {
	color.Yellow("压缩目录: %s", srcDir)
	
	// 创建目标目录
	if err := os.MkdirAll(distDir, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 遍历源目录
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查是否为图片文件
		if !isImageFile(path) {
			return nil
		}

		// 计算相对路径
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// 目标文件路径
		distPath := filepath.Join(distDir, relPath)
		distDirPath := filepath.Dir(distPath)
		
		// 创建目标子目录
		if err := os.MkdirAll(distDirPath, 0755); err != nil {
			return err
		}

		// 压缩文件（保持原扩展名）
		color.Cyan("压缩: %s", relPath)
		if err := compressImageFile(path, distPath, rate); err != nil {
			color.Red("压缩失败: %s - %v", relPath, err)
			return nil // 继续处理其他文件
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("遍历目录失败: %v", err)
	}

	color.Green("✓ 目录压缩完成: %s", distDir)
	return nil
}

// compressFile 压缩单个文件
func compressFile(srcFile, distFile string, rate float64) error {
	color.Yellow("压缩文件: %s", filepath.Base(srcFile))
	
	// 检查是否为图片文件
	if !isImageFile(srcFile) {
		return fmt.Errorf("不支持的文件格式: %s", filepath.Ext(srcFile))
	}

	// 压缩文件
	if err := compressImageFile(srcFile, distFile, rate); err != nil {
		return fmt.Errorf("压缩文件失败: %v", err)
	}

	color.Green("✓ 文件压缩完成: %s", distFile)
	return nil
}

// isImageFile 检查是否为图片文件
func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	
	for _, imgExt := range imageExts {
		if ext == imgExt {
			return true
		}
	}
	return false
}

// compressImageFile 压缩单个图片文件
func compressImageFile(srcFile, distFile string, rate float64) error {
	// 读取源图片
	srcData, err := os.ReadFile(srcFile)
	if err != nil {
		return fmt.Errorf("读取源文件失败: %v", err)
	}

	// 解码图片
	img, format, err := decodeImage(srcData)
	if err != nil {
		// 解码失败，直接复制文件并保持原扩展名
		color.Yellow("⚠️  无法解码图片: %s，直接复制文件", filepath.Base(srcFile))
		if err := os.WriteFile(distFile, srcData, 0644); err != nil {
			return fmt.Errorf("复制文件失败: %v", err)
		}
		
		// 显示复制信息
		color.Green("✓ 文件复制完成: %s", filepath.Base(srcFile))
		color.Cyan("  文件大小: %d bytes", len(srcData))
		return nil
	}

	// 获取原始尺寸
	originalBounds := img.Bounds()
	originalWidth := originalBounds.Dx()
	originalHeight := originalBounds.Dy()

	// 计算新尺寸（按比率缩小）
	newWidth := int(float64(originalWidth) * rate)
	newHeight := int(float64(originalHeight) * rate)

	// 如果新尺寸太小，保持最小尺寸
	if newWidth < 50 {
		newWidth = 50
	}
	if newHeight < 50 {
		newHeight = 50
	}

	// 创建目标文件
	distFileHandle, err := os.Create(distFile)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer distFileHandle.Close()

	// 根据格式进行压缩
	switch format {
	case "jpeg", "jpg":
		// JPEG 压缩：先优化质量，再调整尺寸
		err = compressJPEG(img, distFileHandle, newWidth, newHeight, rate)
	case "png":
		// PNG 压缩：先优化质量，再调整尺寸
		err = compressPNG(img, distFileHandle, newWidth, newHeight, rate)
	case "gif":
		// GIF 压缩：先优化质量，再调整尺寸
		err = compressGIF(img, distFileHandle, newWidth, newHeight, rate)
	default:
		// 不支持的格式，直接复制
		color.Yellow("⚠️  不支持的格式: %s，直接复制文件", format)
		_, err = distFileHandle.Write(srcData)
	}

	if err != nil {
		return fmt.Errorf("压缩图片失败: %v", err)
	}

	// 获取压缩后的文件大小
	distInfo, err := distFileHandle.Stat()
	if err != nil {
		return fmt.Errorf("获取目标文件信息失败: %v", err)
	}

	originalSize := len(srcData)
	compressedSize := distInfo.Size()
	compressionRatio := float64(compressedSize) / float64(originalSize)

	color.Green("✓ 压缩完成: %s", filepath.Base(srcFile))
	color.Cyan("  原始尺寸: %dx%d", originalWidth, originalHeight)
	color.Cyan("  压缩尺寸: %dx%d", newWidth, newHeight)
	color.Cyan("  原始大小: %d bytes", originalSize)
	color.Cyan("  压缩大小: %d bytes", compressedSize)
	color.Cyan("  压缩比率: %.2f%%", compressionRatio*100)

	return nil
}

// decodeImage 解码图片数据
func decodeImage(data []byte) (image.Image, string, error) {
	img, format, err := image.Decode(strings.NewReader(string(data)))
	if err != nil {
		return nil, "", fmt.Errorf("无法解码图片格式: %v", err)
	}
	return img, format, nil
}

// compressJPEG 压缩JPEG图片
func compressJPEG(img image.Image, file *os.File, width, height int, rate float64) error {
	// 计算JPEG质量（基于压缩比率，但优先保证质量）
	quality := int(85 + (rate-0.5)*30) // 质量范围：70-100
	if quality < 70 {
		quality = 70
	}
	if quality > 100 {
		quality = 100
	}

	// 调整图片尺寸
	resizedImg := resizeImage(img, width, height)

	// 编码为JPEG
	return jpeg.Encode(file, resizedImg, &jpeg.Options{Quality: quality})
}

// compressPNG 压缩PNG图片
func compressPNG(img image.Image, file *os.File, width, height int, rate float64) error {
	// 调整图片尺寸
	resizedImg := resizeImage(img, width, height)

	// PNG使用默认压缩（PNG是无损格式，主要通过尺寸调整来减小文件大小）
	return png.Encode(file, resizedImg)
}

// compressGIF 压缩GIF图片
func compressGIF(img image.Image, file *os.File, width, height int, rate float64) error {
	// 调整图片尺寸
	resizedImg := resizeImage(img, width, height)

	// GIF使用默认编码（GIF压缩主要通过尺寸调整）
	return gif.Encode(file, resizedImg, nil)
}

// resizeImage 调整图片尺寸（简化实现）
func resizeImage(img image.Image, width, height int) image.Image {
	// 这里是一个简化的实现
	// 在实际项目中，你可能需要使用专门的图片处理库来进行高质量的缩放
	
	// 获取原始尺寸
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	// 如果尺寸相同，直接返回原图
	if originalWidth == width && originalHeight == height {
		return img
	}

	// 简化的缩放实现：取平均值
	// 注意：这是一个非常简化的实现，实际应该使用双线性插值等算法
	resized := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// 计算缩放比例
	scaleX := float64(originalWidth) / float64(width)
	scaleY := float64(originalHeight) / float64(height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 计算原图对应位置
			srcX := int(float64(x) * scaleX)
			srcY := int(float64(y) * scaleY)
			
			// 确保不越界
			if srcX >= originalWidth {
				srcX = originalWidth - 1
			}
			if srcY >= originalHeight {
				srcY = originalHeight - 1
			}

			// 复制像素
			resized.Set(x, y, img.At(srcX, srcY))
		}
	}

	return resized
} 