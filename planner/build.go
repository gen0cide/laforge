package planner

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gen0cide/laforge/ent"
)

func BuildPlan(ctx context.Context, client *ent.Client) error {
	nodes, err := client.Plan.Query().All(ctx)

	if err != nil {
		log.Fatalf("Failed to Query Plan Nodes %v. Err: %v", nodes, err)
		return err
	}

	var wg sync.WaitGroup

	for _, node := range nodes {
		status, err := node.PlanToStatus(ctx)

		if err != nil {
			log.Fatalf("Failed to Query Status %v. Err: %v", node, err)
			return err
		}

		wg.Add(1)

		go func(wg *sync.WaitGroup, status *ent.Status) {
			defer wg.Done()
			ctx := context.Background()
			defer ctx.Done()
			status.Update().SetState("AWAITING").Save(ctx)
		}(&wg, status)
	}

	wg.Wait()

	for _, node := range nodes {
		wg.Add(1)
		go buildRoutine(node, ctx)
	}

	wg.Wait()

	return nil
}

func buildRoutine(node *ent.Plan, ctx context.Context) {
	prevNodes, err := node.PrevPlan(ctx)

	if err != nil {
		log.Fatalf("Failed to Query Plan Start %v. Err: %v", prevNodes, err)
	}

	for _, prevNode := range prevNodes {
		for {
			prevStatus, err := prevNode.PlanToStatus(ctx)

			if err != nil {
				log.Fatalf("Failed to Query Status %v. Err: %v", prevNode, err)
			}

			if prevNode == nil && prevStatus.Completed {
				break
			}

			time.Sleep(time.Second)
		}
	}

	// go build(node)
	status, err := node.PlanToStatus(ctx)

	if err != nil {
		log.Fatalf("Failed to Query Status %v. Err: %v", node, err)
	}

	status.Update().SetState("COMPLETED").Save(ctx)
	status.Update().SetCompleted(true).Save(ctx)
}
