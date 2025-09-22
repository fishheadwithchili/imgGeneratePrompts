package services_test

import (
	"imgGeneratePrompts/models"
	"testing"

	"github.com/stretchr/testify/suite"
)

// TagServiceTestSuite 是 TagService 的测试套件
// 我们复用 PromptServiceTestSuite 的设置，因为它已经包含了所有需要的东西
type TagServiceTestSuite struct {
	PromptServiceTestSuite
}

// TestCreateAndGetTag 测试创建和获取标签
func (s *TagServiceTestSuite) TestCreateAndGetTag() {
	// 1. 测试创建
	req := &models.CreateTagRequest{Name: "测试标签"}
	tag, err := s.tagSvc.CreateTag(req)
	s.NoError(err)
	s.NotNil(tag)
	s.Equal("测试标签", tag.Name)

	// 2. 测试通过 ID 获取
	fetchedTag, err := s.tagSvc.GetTagByID(tag.ID)
	s.NoError(err)
	s.NotNil(fetchedTag)
	s.Equal(tag.ID, fetchedTag.ID)
	s.Equal("测试标签", fetchedTag.Name)

	// 3. 测试通过名称获取
	fetchedTagByName, err := s.tagSvc.GetTagByName("测试标签")
	s.NoError(err)
	s.NotNil(fetchedTagByName)
	s.Equal(tag.ID, fetchedTagByName.ID)
}

// TestCreateDuplicateTag 测试创建重复标签
func (s *TagServiceTestSuite) TestCreateDuplicateTag() {
	// 第一次创建
	req := &models.CreateTagRequest{Name: "重复标签"}
	tag1, err := s.tagSvc.CreateTag(req)
	s.NoError(err)

	// 第二次创建同名标签
	tag2, err := s.tagSvc.CreateTag(req)
	s.NoError(err)

	// 应该返回已存在的标签，ID 相同
	s.Equal(tag1.ID, tag2.ID)
}

// TestDeleteTag 测试删除标签
func (s *TagServiceTestSuite) TestDeleteTag() {
	// 1. 创建一个未被使用的标签并删除
	tag, _ := s.tagSvc.CreateTag(&models.CreateTagRequest{Name: "待删除标签"})
	err := s.tagSvc.DeleteTag(tag.ID)
	s.NoError(err)

	// 验证它已被删除
	_, err = s.tagSvc.GetTagByID(tag.ID)
	s.Error(err)

	// 2. 创建一个被使用的标签，删除时应该报错
	tagUsed, _ := s.tagSvc.CreateTag(&models.CreateTagRequest{Name: "被使用的标签"})
	s.service.CreatePromptWithImages(&models.CreatePromptRequest{
		PromptText: "关联提示词",
		TagNames:   []string{"被使用的标签"},
	})

	err = s.tagSvc.DeleteTag(tagUsed.ID)
	s.Error(err, "删除被使用的标签时应该返回错误")
	s.Contains(err.Error(), "无法删除标签，还有")
}

// TestGetOrCreateTags 测试批量获取或创建标签
func (s *TagServiceTestSuite) TestGetOrCreateTags() {
	// 准备一个已存在的标签
	s.tagSvc.CreateTag(&models.CreateTagRequest{Name: "已存在"})

	tagNames := []string{"已存在", "新标签1", "新标签2", "  带空格的标签  "}
	tags, err := s.tagSvc.GetOrCreateTags(tagNames)
	s.NoError(err)
	s.Len(tags, 4)

	// 验证数据库中确实有这些标签
	var count int64
	s.db.Model(&models.Tag{}).Where("name IN ?", []string{"已存在", "新标签1", "新标签2", "带空格的标签"}).Count(&count)
	s.Equal(int64(4), count)
}

// TestTagStats 测试标签统计
func (s *TagServiceTestSuite) TestTagStats() {
	// 准备数据
	s.tagSvc.CreateTag(&models.CreateTagRequest{Name: "热门标签"})
	s.tagSvc.CreateTag(&models.CreateTagRequest{Name: "冷门标签"})
	s.service.CreatePromptWithImages(&models.CreatePromptRequest{PromptText: "p1", TagNames: []string{"热门标签"}})
	s.service.CreatePromptWithImages(&models.CreatePromptRequest{PromptText: "p2", TagNames: []string{"热门标签"}})

	stats, err := s.tagSvc.GetTagStats()
	s.NoError(err)

	s.Equal(int64(2), stats["total_tags"])

	// 既然源代码返回了正确的类型，我们就可以直接进行类型断言
	popularTags, ok := stats["popular_tags"].([]map[string]interface{})
	s.True(ok, "popular_tags 应该是 []map[string]interface{} 类型")

	s.Len(popularTags, 1)
	s.Equal("热门标签", popularTags[0]["tag_name"])
	s.Equal(int64(2), popularTags[0]["use_count"])
}

// TestTagService runs the test suite for the tag service
func TestTagService(t *testing.T) {
	// 复用 PromptServiceTestSuite 的 setup 和 teardown 逻辑
	suite.Run(t, new(TagServiceTestSuite))
}
