// ==========================================================================
// 云捷GO自动生成业务逻辑层相关代码，只生成一次，按需修改,再次生成不会覆盖.
// 生成日期：2020-02-18 15:44:13
// 生成路径: app/service/module/job/job_service.go
// 生成人：yunjie
// ==========================================================================
package job

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	jobModel "yj-app/app/model/monitor/job"
	userService "yj-app/app/service/system/user"
	"yj-app/app/service/utils/convert"
	"yj-app/app/service/utils/excel"
	"yj-app/app/service/utils/page"
)

//根据主键查询数据
func SelectRecordById(id int64) (*jobModel.Entity, error) {
	return jobModel.FindOne("job_id", id)
}

//根据主键删除数据
func DeleteRecordById(id int64) bool {
	result, err := jobModel.Delete("job_id", id)
	if err == nil {
		affected, _ := result.RowsAffected()
		if affected > 0 {
			return true
		}
	}

	return false
}

//批量删除数据记录
func DeleteRecordByIds(ids string) int64 {
	idarr := convert.ToInt64Array(ids, ",")
	result, err := jobModel.Delete("job_id in (?)", idarr)
	if err != nil {
		return 0
	}

	nums, _ := result.RowsAffected()

	return nums
}

//添加数据
func AddSave(req *jobModel.AddReq, session *ghttp.Session) (int64, error) {
	//检查任务名称是否存在
	rs := gcron.Search(req.JobName)

	if rs != nil {
		return 0, gerror.New("任务名称已经存在")
	}

	if req.MisfirePolicy == "1" {
		task, err := gcron.Add(req.CronExpression, func() { glog.Println("Every hour on the half hour") }, req.JobName)
		if err != nil || task == nil {
			return 0, err
		}
	} else {
		task, err := gcron.AddOnce(req.CronExpression, func() { glog.Println("Every hour on the half hour") }, req.JobName)
		if err != nil || task == nil {
			return 0, err
		}
	}

	var entity jobModel.Entity
	entity.JobName = req.JobName
	entity.JobGroup = req.JobGroup
	entity.InvokeTarget = req.InvokeTarget
	entity.CronExpression = req.CronExpression
	entity.MisfirePolicy = req.MisfirePolicy
	entity.Concurrent = req.Concurrent
	entity.Status = req.Status
	entity.Remark = req.Remark
	entity.CreateTime = gtime.Now()
	entity.CreateBy = ""

	user := userService.GetProfile(session)

	if user != nil {
		entity.CreateBy = user.LoginName
	}

	result, err := entity.Insert()
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil || id <= 0 {
		return 0, err
	}
	return id, nil
}

//修改数据
func EditSave(req *jobModel.EditReq, session *ghttp.Session) (int64, error) {
	//检查任务名称是否存在
	tmp := gcron.Search(req.JobName)

	if tmp != nil {
		gcron.Remove(req.JobName)
	}

	if req.MisfirePolicy == "1" {
		task, err := gcron.Add(req.CronExpression, func() { glog.Println("Every hour on the half hour") }, req.JobName)
		if err != nil || task == nil {
			return 0, err
		}
	} else {
		task, err := gcron.AddOnce(req.CronExpression, func() { glog.Println("Every hour on the half hour") }, req.JobName)
		if err != nil || task == nil {
			return 0, err
		}
	}

	entity, err := jobModel.FindOne("job_id=?", req.JobId)

	if err != nil {
		return 0, err
	}

	if entity == nil {
		return 0, gerror.New("数据不存在")
	}

	entity.InvokeTarget = req.InvokeTarget
	entity.CronExpression = req.CronExpression
	entity.MisfirePolicy = req.MisfirePolicy
	entity.Concurrent = req.Concurrent
	entity.Status = req.Status
	entity.Remark = req.Remark
	entity.UpdateTime = gtime.Now()
	entity.UpdateBy = ""

	user := userService.GetProfile(session)

	if user == nil {
		entity.UpdateBy = user.LoginName
	}

	result, err := entity.Update()

	if err != nil {
		return 0, err
	}

	rs, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rs, nil
}

//根据条件查询数据
func SelectListAll(params *jobModel.SelectPageReq) (*[]jobModel.Entity, error) {
	return jobModel.SelectListAll(params)
}

//根据条件分页查询数据
func SelectListByPage(params *jobModel.SelectPageReq) (*[]jobModel.Entity, *page.Paging, error) {
	return jobModel.SelectListByPage(params)
}

// 导出excel
func Export(param *jobModel.SelectPageReq) (string, error) {
	result, err := jobModel.SelectListExport(param)
	if err != nil {
		return "", err
	}

	head := []string{"任务ID", "任务名称", "任务组名", "调用目标字符串", "cron执行表达式", "计划执行错误策略（1立即执行 2执行一次 3放弃执行）", "是否并发执行（0允许 1禁止）", "状态（0正常 1暂停）", "创建者", "创建时间", "更新者", "更新时间", "备注信息"}

	url, err := excel.DownlaodExcel(head, result)

	if err != nil {
		return "", err
	}

	return url, nil
}
