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

type Category struct {
	Id              int64 `orm:"pk;auto"`
	Title           string
	Created         time.Time `orm:"index;null"`
	Views           int64     `orm:"index;null"`
	TopicTime       time.Time `orm:"index;null"`
	TopicCount      int64     `orm:"null"`
	TopicLastUserId int64     `orm:"null"`
}

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

type Comment struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}

// 注册数据库
func RegisterDB() {
	if !IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	orm.RegisterModel(new(Category), new(Topic), new(Comment))
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
}

// 判断数据库是否存在
func IsExist(fn string) bool {
	f, err := os.Open(fn)
	defer f.Close()
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// 添加文章分类
func AddCategory(name string) error {
	o := orm.NewOrm()
	cate := &Category{Title: name}
	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate)
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
func DeleteCategory(cid string) error {
	cidNum, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Category{Id: cidNum}
	_, err = o.Delete(cate)
	return err
}

// 获取指定分类
func GetCategory(cid string) (*Category, error) {
	cidNum, err := strconv.ParseInt(cid, 10, 64)
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
	return cate, err
}

// 获取全部文章分类
func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()
	cates := make([]*Category, 0)
	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

// 分类文章数增减
func TopicCountUp(cid string, u bool) error {
	cidNum, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("id", cidNum).One(cate)
	if err != nil {
		return err
	}
	if u {
		cate.TopicCount++
	} else {
		cate.TopicCount--
	}
	_, err = o.Update(cate)
	return err
}

// 添加文章
func AddTopic(title, category, content string) error {
	o := orm.NewOrm()
	topic := &Topic{
		Title:    title,
		Category: category,
		Content:  content,
		Created:  time.Now(),
		Updated:  time.Now(),
	}
	_, err := o.Insert(topic)
	return err
}

// 修改指定文章
func ModifyTopic(tid, title, category, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	err = o.Read(topic)
	if err != nil {
		return err
	}
	topic.Title = title
	topic.Category = category
	topic.Content = content
	topic.Updated = time.Now()
	o.Update(topic)
	return nil
}

// 删除指定文章
func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	_, err = o.Delete(topic)
	return err
}

// 获取指定文章
func GetTopic(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
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
	topic.Views++
	_, err = o.Update(topic)
	return topic, err
}

// 获取全部文章
func GetAllTopics(isDesc bool) ([]*Topic, error) {
	var err error
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")
	if isDesc {
		_, err = qs.OrderBy("-created").All(&topics)
	} else {
		_, err = qs.All(&topics)
	}
	return topics, err
}

// 通过文章ID获取分类
func GetCategoryByTopicId(tid string) (*Category, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
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
	cate := new(Category)
	qs = o.QueryTable("category")
	err = qs.Filter("title", topic.Category).One(cate)
	if err != nil {
		return nil, err
	}
	return cate, nil
}

// 添加评论
func AddReply(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	comment := &Comment{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: time.Now(),
	}
	_, err = o.Insert(comment)
	if err == nil {
		topic := new(Topic)
		qs := o.QueryTable("topic")
		err = qs.Filter("id", tidNum).One(topic)
		if err != nil {
			return err
		}
		topic.ReplyCount++
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}
	return err
}

// 删除指定评论
func DeleteReply(tid, rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	comment := &Comment{Id: ridNum}
	_, err = o.Delete(comment)
	if err == nil {
		topic := new(Topic)
		qs := o.QueryTable("topic")
		err = qs.Filter("id", tidNum).One(topic)
		if err != nil {
			return err
		}
		topic.ReplyCount--
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}
	return err
}

// 获取全部评论
func GetAllReplies(tid string) ([]*Comment, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	replies := make([]*Comment, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).All(&replies)
	return replies, err
}
