// TestStoreEmail 测试邮件存储功能
// 包含以下测试用例：
// 1. 成功存储邮件并验证文件创建
// 2. 测试无效输入情况：
//   - 空发件人地址
//   - 无收件人
//   - 空邮件内容
//
// 测试完成后会自动清理临时文件
package mail_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/YoPost/internal/config"
	"github.com/YoPost/internal/mail"
	"github.com/stretchr/testify/assert"
)

func TestStoreEmail(t *testing.T) {
	// 初始化测试配置
	cfg := &config.Config{}
	core, err := mail.NewCore(cfg)
	assert.NoError(t, err)

	// 测试成功存储
	t.Run("successful storage", func(t *testing.T) {
		err := core.StoreEmail("test@example.com", []string{"recipient@example.com"}, "Test email content")
		assert.NoError(t, err)

		// 验证文件是否创建
		tmpDir := os.TempDir()
		matches, err := filepath.Glob(filepath.Join(tmpDir, "yopost_emails", "*.eml"))
		assert.NoError(t, err)
		assert.NotEmpty(t, matches)
	})

	// 测试无效输入
	t.Run("empty from address", func(t *testing.T) {
		err := core.StoreEmail("", []string{"recipient@example.com"}, "content")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty from address")
	})

	t.Run("no recipients", func(t *testing.T) {
		err := core.StoreEmail("test@example.com", []string{}, "content")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no recipients specified")
	})

	t.Run("empty content", func(t *testing.T) {
		err := core.StoreEmail("test@example.com", []string{"recipient@example.com"}, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty email content")
	})

	// 清理测试文件
	t.Cleanup(func() {
		os.RemoveAll(filepath.Join(os.TempDir(), "yopost_emails"))
	})
}
