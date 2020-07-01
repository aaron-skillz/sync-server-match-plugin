package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/aaron-skillz/sync-server-go/runtime"
	"strconv"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	if err := initializer.RegisterRpc("create_skillz_match", createMatchRPC); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}
	if err := initializer.RegisterMatch("skillz-match", func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return &Match{}, nil
	}); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}
	logger.Info("Skillz Match Plugin Loaded")
	return nil
}

func createMatchRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	params := make(map[string]interface{})
	if err := json.Unmarshal([]byte(payload), &params); err != nil {
		return "", err
	}

	modulename := "skillz-match" // Name with which match handler was registered in InitModule, see example above.
	if matchId, err := nk.MatchCreate(ctx, modulename, params); err != nil {
		return "", err
	} else {
		return matchId, nil
	}
}

type MatchState struct {
	presences map[string]runtime.Presence
}

type Match struct{}

func (m *Match) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	state := &MatchState{
		presences: make(map[string]runtime.Presence),
	}
	tickRate := 1
	label := ""
	return state, tickRate, label
}

func (m *Match) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	acceptUser := true
	return state, acceptUser, ""
}

func (m *Match) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	mState, _ := state.(*MatchState)
	for _, p := range presences {
		mState.presences[p.GetUserId()] = p
	}

	return mState
}

func (m *Match) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	mState, _ := state.(*MatchState)
	for _, p := range presences {
		delete(mState.presences, p.GetUserId())
	}

	return mState
}

func (m *Match) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	mState, _ := state.(*MatchState)
	for _, presence := range mState.presences {
		logger.Info("Presence %v named %v", presence.GetUserId(), presence.GetUsername())
	}

	for _, message := range messages {
		logger.Info("Received %v from %v", string(message.GetData()), message.GetUserId())

		dispatcher.BroadcastMessage(1, message.GetData(), []runtime.Presence{message}, nil, true)
	}

	return mState
}

func (m *Match) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	message := "Server shutting down in " + strconv.Itoa(graceSeconds) + " seconds."
	dispatcher.BroadcastMessage(2, []byte(message), nil, nil, true)
	return state
}
