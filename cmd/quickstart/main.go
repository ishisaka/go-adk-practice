// ADK サンプル
// https://google.github.io/adk-docs/get-started/go/#set-your-api-key
package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	// モデルのインスタンスを作る
	model, err := gemini.NewModel(ctx, "gemini-3-flash-preview", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	// エージェントのインスタンスを作る
	timeAgent, err := llmagent.New(llmagent.Config{
		// エージェントの名前を設定する
		Name: "hello_time_agent",
		// エージェントが使用するモデルを設定する
		Model: model,
		// エージェントの説明を設定する
		Description: "Tells the current time in a specified city.",
		// エージェントへの指示を設定する
		Instruction: "You are a helpful assistant that tells the current time in a city.",
		// ツールを追加する(Web検索ツールを追加)
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// ランチャーの設定を作る
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(timeAgent),
	}

	// ランチャーを作って実行する
	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
