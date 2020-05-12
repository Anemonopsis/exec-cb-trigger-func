package exec_cb_trigger_func

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/cloudbuild/v1"
	"log"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type CloudBuildExecuteTriggerCommandArgs struct {
	GcpProjectID string `json:"gcp_project_id"`
	TriggerId    string `json:"trigger_id"`
	Branch       string `json:"branch"`
	SHA          string `json:"sha"`
	Tag          string `json:"tag"`
}

type CloudBuildTrigger struct {
	ID         string
	RepoSource *cloudbuild.RepoSource
}

type CloudBuildService struct {
	cbs *CloudBuildService
}

func NewCloudBuildService(ctx context.Context) *CloudBuildService {
	// TODO
	return &CloudBuildService{cbs: NewCloudBuildService(ctx)}
}

func NewCloudBuildExecuteTriggerCommandArgs(data []byte) (*CloudBuildExecuteTriggerCommandArgs, error) {
	var args CloudBuildExecuteTriggerCommandArgs
	if err := json.Unmarshal(data, &args); err != nil {
		return nil, fmt.Errorf("pub sub message data unmershall failed. check your data if data format is correct json")
	}

	// TODO
	//if args.GcpProjectID == "" || args.InstanceName == "" {
	//	return nil, fmt.Errorf("gcp_project_id, instance_name are required. please including they are in pub sub message data")
	//}

	return &args, nil
}

// TODO 処理確定させてメソッドに分割する
func ExecCloudBuildTriggerFunc(ctx context.Context, m PubSubMessage) error {
	_, err := NewCloudBuildExecuteTriggerCommandArgs(m.Data)
	if err != nil {
		return err
	}

	// Service初期化
	ser, err := cloudbuild.NewService(ctx)
	if err != nil {
		return err
	}
	tSer := cloudbuild.NewProjectsTriggersService(ser)

	// Trigger一覧取得
	listCall := tSer.List("project_id")
	list, err := listCall.Do()
	if err != nil {
		return err
	}

	// TriggerからSearch
	var targetTrigger *CloudBuildTrigger
	for _, trigger := range list.Triggers {
		if !IsTargetTrigger(trigger) {
			continue
		}
		targetTrigger = &CloudBuildTrigger{
			ID:         trigger.Id,
			RepoSource: trigger.Build.Source.RepoSource,
		}
	}
	if targetTrigger == nil {
		return fmt.Errorf("target trigger is not found")
	}

	// Trigger Run
	op, err := tSer.Run("", targetTrigger.ID, targetTrigger.RepoSource).Do()
	if err != nil {
		return err
	}

	log.Println(op)
	return nil
}

func IsTargetTrigger(trigger *cloudbuild.BuildTrigger) bool {
	// TODO
	return true
}
