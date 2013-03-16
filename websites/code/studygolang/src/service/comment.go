package service

import (
	"logger"
	"model"
	"strconv"
	"time"
	"util"
)

// 获得某人在某种类型最近的评论
func FindRecentComments(uid, objtype int) []*model.Comment {
	comments, err := model.NewComment().Where("uid=" + strconv.Itoa(uid) + " AND objtype=" + strconv.Itoa(objtype)).Order("ctime DESC").Limit("0, 5").FindAll()
	if err != nil {
		logger.Errorln("comment service FindRecentComments error:", err)
		return nil
	}
	return comments
}

// 某类型的评论总数
func CommentsTotal(objtype int) (total int) {
	total, err := model.NewComment().Where("objtype=" + strconv.Itoa(objtype)).Count()
	if err != nil {
		logger.Errorln("comment service CommentsTotal error:", err)
		return
	}
	return
}

var commenters = make(map[string]Commenter)

// 评论接口
type Commenter interface {
	// 评论回调接口，用于更新对象自身需要更新的数据
	UpdateComment(int, int, int, string)
}

// 注册评论对象，使得某种类型（帖子、博客等）可以被评论
func RegisterCommentObject(objname string, commenter Commenter) {
	if commenter == nil {
		panic("service: Register commenter is nil")
	}
	if _, dup := commenters[objname]; dup {
		panic("service: Register called twice for commenter " + objname)
	}
	commenters[objname] = commenter
}

// 发表评论。入topics_reply库，更新topics和topics_ex库
// objname 注册的评论对象名
func PostComment(objid, objtype, uid int, content string, objname string) error {
	comment := model.NewComment()
	comment.Objid = objid
	comment.Objtype = objtype
	comment.Uid = uid
	comment.Content = content

	// TODO:评论楼层怎么处理，避免冲突？最后的楼层信息保存在内存中？

	// 暂时只是从数据库中取出最后的评论楼层
	stringBuilder := util.NewBuffer()
	stringBuilder.Append("objid=").AppendInt(objid).Append(" AND objtype=").AppendInt(objtype)
	tmpCmt, err := model.NewComment().Where(stringBuilder.String()).Order("ctime DESC").Find()
	if err != nil {
		logger.Errorln("post comment service error:", err)
		return err
	} else {
		comment.Floor = tmpCmt.Floor + 1
	}
	// 入评论库
	cid, err := comment.Insert()
	if err != nil {
		logger.Errorln("post comment service error:", err)
		return err
	}
	// 回调，不关心处理结果（有些对象可能不需要回调）
	if commenter, ok := commenters[objname]; ok {
		logger.Debugf("评论[objid:%d] [objtype:%d] [uid:%d] 成功，通知被评论者更新", objid, objtype, uid)
		go commenter.UpdateComment(cid, objid, uid, time.Now().Format("2006-01-02 15:04:05"))
	}

	return nil
}
