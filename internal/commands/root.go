package commands

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/your-repo/cyben-zen-tools/internal/config"
)

// NewRootCommand 创建根命令
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cyber-zen",
		Short: "Cyber Zen Tools - 跨平台命令行工具集",
		Long: `Cyber Zen Tools 是一个跨平台的命令行工具集，
提供常用的开发工具和快捷命令。

支持的命令:
  gcm        - Git 提交并推送
  status     - 显示工具状态
  uninstall  - 卸载程序
  compress   - 压缩图片文件`,
		Version: "1.0.0",
		Run: func(cmd *cobra.Command, args []string) {
			// 如果没有子命令，显示帮助
			if len(args) == 0 {
				_ = cmd.Help()
			}
		},
	}

	// 添加子命令
	rootCmd.AddCommand(newGcmCommand())
	rootCmd.AddCommand(newStatusCommand())
	rootCmd.AddCommand(newUninstallCommand())
	rootCmd.AddCommand(newCompressCommand())
	


	return rootCmd
}



// newGcmCommand 创建 gcm 命令
func newGcmCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gcm [message]",
		Short: "Git 提交并推送",
		Long: `快速执行 Git 提交和推送操作：
  1. git add .
  2. git commit -m "message" --no-verify
  3. git push

如果没有提供提交信息，默认使用 "update"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGcm(args)
		},
	}

	return cmd
}



// newStatusCommand 创建状态命令
func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "显示运行时状态",
		Run: func(cmd *cobra.Command, args []string) {
			showStatus()
		},
	}

	return cmd
}

// newUninstallCommand 创建卸载命令
func newUninstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "卸载程序",
		Long: `从系统中卸载 Cyber Zen Tools。

卸载操作会：
1. 删除 /usr/local/bin/cyber-zen 文件
2. 清理构建目录（如果存在）

注意：卸载需要管理员权限`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUninstall()
		},
	}

	return cmd
}





// runUninstall 执行卸载操作
func runUninstall() error {
	color.Yellow("开始卸载 Cyber Zen Tools...")
	
	// 检查是否已安装
	installPath := "/usr/local/bin/cyber-zen"
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		color.Yellow("程序未安装: %s", installPath)
		return nil
	}
	
	// 删除安装文件
	color.Yellow("删除安装文件: %s", installPath)
	cmd := exec.Command("sudo", "rm", "-f", installPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("删除安装文件失败: %v", err)
	}
	
	// 清理构建目录
	color.Yellow("清理构建目录...")
	if err := os.RemoveAll("build"); err != nil {
		color.Yellow("清理构建目录失败: %v", err)
	}
	
	color.Green("✓ 卸载完成！")
	return nil
}



// runGcm 执行 Git 提交和推送
func runGcm(args []string) error {
	// 获取提交信息
	msg := "update"
	if len(args) > 0 {
		msg = args[0]
	}

	color.Green("开始执行 Git 操作...")
	color.Cyan("提交信息: %s", msg)

	// 检查是否在 Git 仓库中
	if err := checkGitRepo(); err != nil {
		return err
	}

	// 执行 git add .
	color.Yellow("执行: git add .")
	if err := execGitCommand("add", "."); err != nil {
		return fmt.Errorf("git add 失败: %v", err)
	}
	color.Green("✓ git add . 完成")

	// 执行 git commit
	color.Yellow("执行: git commit -m \"%s\" --no-verify", msg)
	if err := execGitCommand("commit", "-m", msg, "--no-verify"); err != nil {
		fmt.Errorf("git commit 失败: %v", err)
	}
	color.Green("✓ git commit 完成")

	// 执行 git push
	color.Yellow("执行: git push")
	if err := execGitCommand("push"); err != nil {
		return fmt.Errorf("git push 失败: %v", err)
	}
	color.Green("✓ git push 完成")

	color.Green("🎉 Git 操作完成！")
	return nil
}

// checkGitRepo 检查是否在 Git 仓库中
func checkGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("当前目录不是 Git 仓库")
	}
	return nil
}



// execGitCommand 执行 Git 命令
func execGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	return cmd.Run()
}





// showStatus 显示工具状态
func showStatus() {
	color.Green("=== Cyben Zen Tools 状态 ===")
	
	// 显示安装目录
	installDir := config.GetInstallDir()
	color.Cyan("安装目录: %s", installDir)
	
	// 显示版本信息
	color.Cyan("版本: 1.0.0")
	color.Cyan("平台: %s/%s", runtime.GOOS, runtime.GOARCH)
	
	// 检查 Git 是否可用
	if _, err := exec.LookPath("git"); err == nil {
		color.Green("✓ Git 可用")
	} else {
		color.Red("✗ Git 不可用")
	}
	
	// 检查 bash 是否可用
	if _, err := exec.LookPath("bash"); err == nil {
		color.Green("✓ Bash 可用")
	} else {
		color.Red("✗ Bash 不可用")
	}
} 

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

用法:
  cyber-zen compress -src "源文件或文件夹" -dist "目标路径" -rate "压缩比率"

参数:
  -src   源文件或文件夹路径（必需）
  -dist  目标路径（可选，默认当前目录）
  -rate  压缩比率 0.1-1.0（可选，默认0.8）

示例:
  cyber-zen compress -src "images/" -dist "compressed/" -rate 0.7
  cyber-zen compress -src "photo.jpg" -rate 0.5`,
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
		// 压缩单个文件
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
		// 文件
		return filepath.Join(dir, fmt.Sprintf("%s_%s%s", name, timestamp, ext))
	} else {
		// 目录
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

		// 压缩文件
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
		return fmt.Errorf("解码图片失败: %v", err)
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