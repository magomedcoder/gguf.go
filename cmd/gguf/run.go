package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/magomedcoder/gguf.go/runtime"
	"github.com/magomedcoder/gguf.go/sampler"
)

// runRun выполняет генерацию текста
func runRun(args []string) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	modelPath := fs.String("m", "", "путь к файлу GGUF")
	prompt := fs.String("p", "", "текст промпта")
	maxTokens := fs.Int("n", 128, "максимум новых токенов")
	temp := fs.Float64("temp", 0, "температура sampling (0 = greedy)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *modelPath == "" {
		return fmt.Errorf("использование: gguf run -m файл.gguf -p \"промпт\" [-n 128] [--temp 0.7]")
	}
	
	if *prompt == "" {
		return fmt.Errorf("укажите промпт через -p")
	}

	engine, err := runtime.Load(*modelPath)
	if err != nil {
		return err
	}

	ctx, err := engine.NewContext()
	if err != nil {
		return err
	}

	rng := rand.New(rand.NewPCG(0, 0))
	samp := sampler.Temperature(float32(*temp), rng)

	err = ctx.GenerateStream(*prompt, runtime.GenerateParams{
		MaxTokens: *maxTokens,
		Sampler:   samp,
	}, os.Stdout)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout)
	return nil
}
