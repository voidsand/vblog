package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	_DB_NAME        = "/home/yu/workspace/db/vblog.db"
	_SQLITE3_DRIVER = "sqlite3"
)

// 分类结构映射
type Category struct {
	Id              int64 `orm:"pk;auto"`
	Title           string
	Created         time.Time `orm:"index;null"`
	Views           int64     `orm:"index;null"`
	TopicTime       time.Time `orm:"index;null"`
	TopicCount      int64     `orm:"null"`
	TopicLastUserId int64     `orm:"null"`
}

// 文章结构映射
type Topic struct {
	Id              int64
	Uid             int64 `orm:"null"`
	Title           string
	Category        string
	Content         string    `orm:"size(5000)"`
	Attachment      string    `orm:"null"`
	Created         time.Time `orm:"index"`
	Updated         time.Time `orm:"index"`
	Views           int64     `orm:"index;null"`
	Author          string    `orm:"null"`
	ReplyTime       time.Time `orm:"index;null"`
	ReplyCount      int64     `orm:"null"`
	ReplyLastUserId int64     `orm:"null"`
}

// 回复结构映射
type Reply struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}

// 注册数据库
func RegisterDB() {
	// 判断数据库是否存在，不存在则创建
	if !IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}
	// 注册orm的Model，Driver和DataBase
	orm.RegisterModel(new(Category), new(Topic), new(Reply))
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
}

// 判断数据库是否存在
func IsExist(fName string) bool {
	f, err := os.Open(fName)
	defer f.Close()
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// 添加文章分类
func AddCategory(cTitle string) error {
	o := orm.NewOrm()
	cate := &Category{Title: cTitle}
	qs := o.QueryTable("category")
	err := qs.Filter("title", cTitle).One(cate)
	if err == nil {
		return err
	}
	_, err = o.Insert(cate)
	if err != nil {
		return err
	}
	return nil
}

// 删除文章分类
func DeleteCategory(cId string) error {
	cate, err := GetCategory(cId)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	_, err = o.Delete(cate)
	if err != nil {
		return err
	}
	return nil
}

// 通过分类ID获取指定分类
func GetCategory(cId string) (*Category, error) {
	cidNum, err := strconv.ParseInt(cId, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("id", cidNum).One(cate)
	if err != nil {
		return nil, err
	}
	return cate, nil
}

// 获取全部文章分类
func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()
	cates := make([]*Category, 0)
	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	if err != nil {
		return nil, err
	}
	return cates, nil
}

// 添加文章
func AddTopic(tTitle, cId, tContent, tAttachment string) (string, error) {
	cate, err := GetCategory(cId)
	if err != nil {
		return "", err
	}
	topic := &Topic{
		Title:      tTitle,
		Category:   cate.Title,
		Content:    tContent,
		Attachment: tAttachment,
		Created:    time.Now(),
		Updated:    time.Now(),
	}
	o := orm.NewOrm()
	tIdNum, err := o.Insert(topic)
	if err != nil {
		return "", err
	}
	tId := strconv.FormatInt(tIdNum, 10)

	cate.TopicCount, err = GetTopicCountByCategory(cId)
	if err != nil {
		return "", err
	}
	_, err = o.Update(cate)
	if err != nil {
		return "", err
	}
	return tId, err
}

// 修改指定文章
func ModifyTopic(tId, tTitle, cId, TContent, tAttachment string) error {
	// 通过文章ID获取指定分类映射
	oldCate, err := GetCategoryByTopic(tId)
	if err != nil {
		return err
	}
	// 通过分类ID获取指定分类映射
	newCate, err := GetCategory(cId)
	if err != nil {
		return err
	}
	// 通过文章ID获取指定文章映射
	topic, err := GetTopic(tId)
	if err != nil {
		return err
	}
	// 更新指定文章映射内容
	topic.Title = tTitle
	topic.Category = newCate.Title
	topic.Content = TContent
	oldAttachment := topic.Attachment
	topic.Attachment = tAttachment
	topic.Updated = time.Now()
	o := orm.NewOrm()
	_, err = o.Update(topic)
	if err != nil {
		return err
	}
	// 删除旧的附件
	if len(oldAttachment) > 0 {
		os.Remove(path.Join("attachment", tId, oldAttachment))
	}
	// 更新前后分类文章数
	oldCate.TopicCount, err = GetTopicCountByCategory(strconv.FormatInt(oldCate.Id, 10))
	if err != nil {
		return err
	}
	_, err = o.Update(oldCate)
	if err != nil {
		return err
	}
	newCate.TopicCount, err = GetTopicCountByCategory(cId)
	if err != nil {
		return err
	}
	_, err = o.Update(newCate)
	if err != nil {
		return err
	}
	return nil
}

// 删除指定文章
func DeleteTopic(tId string) error {
	cate, err := GetCategoryByTopic(tId)
	if err != nil {
		return err
	}
	topic, err := GetTopic(tId)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	_, err = o.Delete(topic)
	if err != nil {
		return err
	}
	cate.TopicCount, err = GetTopicCountByCategory(strconv.FormatInt(cate.Id, 10))
	if err != nil {
		return err
	}
	_, err = o.Update(cate)
	if err != nil {
		return err
	}
	return nil
}

// 获取指定文章
func GetTopic(tId string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tId, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

// 获取全部文章
func GetAllTopics(cId string, isDesc bool) ([]*Topic, error) {
	topics := make([]*Topic, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("topic")
	if isDesc {
		if len(cId) > 0 {
			cate, err := GetCategory(cId)
			if err != nil {
				return nil, err
			}
			qs = qs.Filter("category", cate.Title)
		}
		_, err := qs.OrderBy("-created").All(&topics)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := qs.All(&topics)
		if err != nil {
			return nil, err
		}
	}
	return topics, nil
}

// 添加评论
func AddReply(tId, rName, rContent string) error {
	topic, err := GetTopic(tId)
	if err != nil {
		return err
	}
	relpy := &Reply{
		Tid:     topic.Id,
		Name:    rName,
		Content: rContent,
		Created: time.Now(),
	}
	o := orm.NewOrm()
	_, err = o.Insert(relpy)
	if err != nil {
		return err
	}
	topic.ReplyCount, err = GetReplyCountByTopic(tId)
	if err != nil {
		return err
	}
	_, err = o.Update(topic)
	if err != nil {
		return err
	}
	return nil
}

// 删除指定评论
func DeleteReply(rId string) error {
	topic, err := GetTopicByReply(rId)
	if err != nil {
		return err
	}
	relpy, err := GetReply(rId)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	_, err = o.Delete(relpy)
	if err != nil {
		return err
	}
	topic.ReplyCount, err = GetReplyCountByTopic(strconv.FormatInt(topic.Id, 10))
	if err != nil {
		return err
	}
	_, err = o.Update(topic)
	if err != nil {
		return err
	}
	return nil
}

// 通过回复ID获取指定回复
func GetReply(rId string) (*Reply, error) {
	rIdNum, err := strconv.ParseInt(rId, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	reply := new(Reply)
	qs := o.QueryTable("reply")
	err = qs.Filter("id", rIdNum).One(reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// 获取全部评论
func GetAllReplies(tId string) ([]*Reply, error) {
	tIdNum, err := strconv.ParseInt(tId, 10, 64)
	if err != nil {
		return nil, err
	}
	replies := make([]*Reply, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("reply")
	_, err = qs.Filter("tid", tIdNum).All(&replies)
	if err != nil {
		return nil, err
	}
	return replies, nil
}

// 通过文章ID获取分类
func GetCategoryByTopic(tId string) (*Category, error) {
	tIdNum, err := strconv.ParseInt(tId, 10, 64)
	if err != nil {
		return nil, err
	}
	topic := new(Topic)
	o := orm.NewOrm()
	qs := o.QueryTable("topic")
	err = qs.Filter("id", tIdNum).One(topic)
	if err != nil {
		return nil, err
	}
	cate := new(Category)
	qs = o.QueryTable("category")
	err = qs.Filter("title", topic.Category).One(cate)
	if err != nil {
		return nil, err
	}
	return cate, nil
}

// 通过回复ID获取文章
func GetTopicByReply(rId string) (*Topic, error) {
	rIdNum, err := strconv.ParseInt(rId, 10, 64)
	if err != nil {
		return nil, err
	}
	reply := new(Reply)
	o := orm.NewOrm()
	qs := o.QueryTable("reply")
	err = qs.Filter("id", rIdNum).One(reply)
	if err != nil {
		return nil, err
	}
	topic := new(Topic)
	qs = o.QueryTable("topic")
	err = qs.Filter("id", reply.Tid).One(topic)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

// 通过分类ID获取分类下文章数
func GetTopicCountByCategory(cId string) (int64, error) {
	cate, err := GetCategory(cId)
	if err != nil {
		return -1, err
	}
	o := orm.NewOrm()
	qs := o.QueryTable("topic")
	tc, err := qs.Filter("category", cate.Title).Count()
	if err != nil {
		return -1, err
	}
	return tc, nil
}

// 通过文章ID获取文章下回复数
func GetReplyCountByTopic(tId string) (int64, error) {
	topic, err := GetTopic(tId)
	if err != nil {
		return -1, err
	}
	o := orm.NewOrm()
	qs := o.QueryTable("reply")
	rc, err := qs.Filter("tid", topic.Id).Count()
	if err != nil {
		return -1, err
	}
	return rc, nil
}

// 修改文章浏览次数
func TopicViewsChange(tId string, up bool) error {
	topic, err := GetTopic(tId)
	if err != nil {
		return err
	}
	if up {
		topic.Views++
	} else {
		topic.Views--
	}
	o := orm.NewOrm()
	_, err = o.Update(topic)
	if err != nil {
		return err
	}
	return nil
}

// 修改分类总浏览次数
func TotalViewsChange(cId string) error {
	var tvc int64
	cate, err := GetCategory(cId)
	if err != nil {
		return err
	}
	topics := make([]*Topic, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("topic")
	_, err = qs.Filter("category", cate.Title).All(&topics)
	if err != nil {
		return err
	}
	for i := range topics {
		tvc += topics[i].Views
	}
	cate.Views = tvc
	_, err = o.Update(cate)
	if err != nil {
		return err
	}
	return nil
}
