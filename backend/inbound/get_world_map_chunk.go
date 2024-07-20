package inbound

import (
	"encoding/json"
	"log/slog"
)

type GetWorldMapChunkReq struct {
	StartX float64 `json:"start_x"`
	StartY float64 `json:"start_y"`
	EndX   float64 `json:"end_x"`
	EndY   float64 `json:"end_y"`
}

type GetWorldMapChunkResp struct {
	Chunk [][]float64 `json:"chunk"`
}

func GetWorldMapChunk(conn connection) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		var r GetWorldMapChunkReq
		if err := json.Unmarshal(req, &r); err != nil {
			return nil, err
		}

		slog.Info("GetWorldMapChunk", "start_x", r.StartX, "start_y", r.StartY, "end_x", r.EndX, "end_y", r.EndY)

		// world, err := world_map.Get(conn.Context().GameMapID)
		// if err != nil {
		// 	return nil, err
		// }

		// chunk, err := world.GetElevationMap(r.StartX, r.StartY, r.EndX, r.EndY)
		// if err != nil {
		// 	return nil, err
		// }

		// slog.Info("GetWorldMapChunk", "chunk", chunk)

		// return GetWorldMapChunkResp{
		// 	Chunk: chunk,
		// }, nil

		return GetWorldMapChunkResp{}, nil
	}
}
