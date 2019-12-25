package dbops

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/pili-video-server/api/def"
	"github.com/pili-video-server/api/utils"
)

func AddUser(userName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO users (username, pwd) VALUES (?, ?)")
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	_, err = stmtIns.Exec(userName, pwd)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func GetUserCredential(userName string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT pwd FROM users WHERE username = ?")
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	var pwd string
	err = stmtOut.QueryRow(userName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	defer stmtOut.Close()
	return pwd, nil
}

func GetUser(userName string) (*def.User, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, pwd FROM users WHERE username = ?")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	var id int
	var pwd string
	err = stmtOut.QueryRow(userName).Scan(&id, &pwd)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	result := &def.User{Id: id, Pwd: pwd, Username: userName}
	defer stmtOut.Close()
	return result, nil
}

func GetUserId(userName string) (int, error) {
	stmtOut, err := dbConn.Prepare("SELECT id FROM users WHERE username = ?")
	if err != nil {
		log.Printf("get user id db prepare error!\n")
		return -1, err
	}

	var id int
	err = stmtOut.QueryRow(userName).Scan(&id)
	if err != nil {
		log.Printf("query row error:%v!\n", err)
		return -1, err
	}

	defer stmtOut.Close()
	return id, nil
}

func ModifyUserInfo(userName string) error {
	//没有过多信息
	return nil
}

func ModifyUserPwd(userName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("UPDATE users SET pwd = ? WHERE username = ?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}
	_, err = stmtIns.Exec(pwd, userName)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func AddNewVideo(aid int, vname string) (*def.VideoInfo, error) {
	vid, _ := utils.NewUUID()

	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05")

	stmtIns, err := dbConn.Prepare(`INSERT INTO video_info
	(id, author_id, name, disply_ctime) VALUES(?, ?, ?, ?)`)
	if err != nil {
		log.Printf("insert db prepare error\n")
		return nil, err
	}

	_, err = stmtIns.Exec(vid, aid, vname, ctime)
	if err != nil {
		return nil, err
	}

	vInfo := &def.VideoInfo{Id: vid, AuthorId: aid, Name: vname, DisplayCtime: ctime}

	defer stmtIns.Close()
	return vInfo, nil
}

func DeleteVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_info WHERE id = ?")
	if err != nil {
		log.Printf("delete db prepare error!\n")
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}

	defer stmtDel.Close()
	return nil
}

func GetVideoInfo(vid string) (*def.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, name, author_id, disply_ctime FROM video_info WHERE id = ?")
	if err != nil {
		log.Printf("get db prepare error!\n")
		return nil, err
	}

	var id string
	var name string
	var authorId int
	var displayCTime string
	err = stmtOut.QueryRow(vid).Scan(&id, &name, &authorId, &displayCTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	res := &def.VideoInfo{Id: id, Name: name, AuthorId: authorId, DisplayCtime: displayCTime}
	defer stmtOut.Close()
	return res, nil
}

func ListVideoInfo(uname string, from, n int) ([]*def.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELECT video_info.id, video_info.author_id, video_info.name, video_info.disply_ctime, modulars.name FROM video_info,modulars,users
	WHERE users.username=(?) AND video_info.modular = modulars.id AND users.id = video_info.author_id
	ORDER BY video_info.create_time DESC LIMIT (?),(?)`)
	if err != nil {
		log.Printf("list db prepare error!\n")
		return nil, err
	}

	var res []*def.VideoInfo
	rows, err := stmtOut.Query(uname, from, n)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var vid string
		var aid int
		var vname string
		var displayTime string
		var modular string
		err := rows.Scan(&vid, &aid, &vname, &displayTime, &modular)
		if err != nil {
			return nil, err
		}
		videoInfo := &def.VideoInfo{Id: vid, AuthorId: aid, Name: vname, DisplayCtime: displayTime, Modular: modular}
		res = append(res, videoInfo)
	}

	defer stmtOut.Close()
	return res, nil
}

func ListVideoInfoMod(mod string, from, n int) ([]*def.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELECT video_info.id, video_info.author_id, video_info.name, video_info.disply_ctime, modulars.name 
	FROM video_info,modulars WHERE modulars.name=?
	ORDER BY video_info.create_time DESC LIMIT ?,?`)
	if err != nil {
		log.Printf("list db prepare error!\n")
		return nil, err
	}

	var res []*def.VideoInfo
	rows, err := stmtOut.Query(mod, from, n)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var vid string
		var aid int
		var vname string
		var displayTime string
		var modular string
		err := rows.Scan(&vid, &aid, &vname, &displayTime, &modular)
		if err != nil {
			return nil, err
		}
		videoInfo := &def.VideoInfo{Id: vid, AuthorId: aid, Name: vname, DisplayCtime: displayTime, Modular: modular}
		res = append(res, videoInfo)
	}

	defer stmtOut.Close()
	return res, nil
}

//点赞与取消点赞
func LikeVideo(vid string, uname string) error {
	conn := Pool.Get()
	if conn == nil {
		log.Printf("redis error!\n")
	}
	defer conn.Close()

	// ustr := string(uid)
	likestr := "like_" + vid
	yes, _ := redis.Bool(conn.Do("sismember", likestr, uname))
	// log.Printf("user:%v the video:%v is like?:%v\n", uname, likestr, yes)
	if !yes {
		_, err = conn.Do("sadd", likestr, uname)
		if err != nil {
			return err
		}
	} else {
		_, err = conn.Do("srem", likestr, uname)
		if err != nil {
			return err
		}
	}

	return nil
}

//获取点赞数
func LikeCount(vid string) (int, error) {
	conn := Pool.Get()
	if conn == nil {
		log.Printf("redis conn error!\n")
	}
	defer conn.Close()

	likestr := "like_" + vid
	res, err := redis.Int(conn.Do("scard", likestr))
	if err != nil {
		return 0, err
	}

	return res, nil

}

//判断用户是否点赞
func IsLike(vid string, uname string) (bool, error) {
	conn := Pool.Get()
	if conn == nil {
		log.Printf("redis conn error!\n")
	}
	defer conn.Close()

	likestr := "like_" + vid
	yes, _ := redis.Bool(conn.Do("sismember", likestr, uname))

	return yes, nil
}

func AddNewComment(aid int, vid, content string) error {
	cid, _ := utils.NewUUID()
	stmtIns, err := dbConn.Prepare("INSERT INTO comments(id, author_id, video_id, comment) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Printf("add new comment db prepare error!\n")
		return err
	}

	_, err = stmtIns.Exec(cid, aid, vid, content)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

func ListComments(vid string, from, n int) ([]*def.CommentInfo, error) {
	var res []*def.CommentInfo
	stmtOut, err := dbConn.Prepare(`SELECT comments.id, video_info.name , users.username, comment 
	FROM comments,video_info,users WHERE users.id=video_info.author_id 
	AND comments.video_id=video_info.id AND video_info.id = ?  
	ORDER BY comments.create_time DESC LIMIT ?,?`)
	if err != nil {
		log.Printf("list comments db prepare error!\n")
		return nil, err
	}

	rows, err := stmtOut.Query(vid, from, n)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var cid, vname, uname, content string
		err := rows.Scan(&cid, &vname, &uname, &content)
		if err != nil {
			return nil, err
		}

		commentInfo := &def.CommentInfo{Id: cid, VideoName: vname, AuthorName: uname, Content: content}
		res = append(res, commentInfo)
	}

	defer stmtOut.Close()
	return res, nil
}

func DeleteComment(cid string) error {
	stmtIns, err := dbConn.Prepare("DELETE FROM comments WHERE id = ?")
	if err != nil {
		log.Printf("delete comment db prepare error!\n")
		return err
	}

	_, err = stmtIns.Exec(cid)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}
