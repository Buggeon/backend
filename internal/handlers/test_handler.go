// Buggeon - SelfHosted service for bug and task tracking
// Copyright (C) 2026 DEVE corp.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package handlers

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/loveholidays/excalidraw-decrypt/pkg/excalidrawdecrypt"
)

type ExcalidrawData struct {
	Type     string    `json:"type"`
	Version  int       `json:"version"`
	Source   string    `json:"source"`
	Elements []Element `json:"elements"`
	AppState AppState  `json:"appState"`
}

type Element struct {
	ID              string        `json:"id"`
	Type            string        `json:"type"`
	X               float64       `json:"x"`
	Y               float64       `json:"y"`
	Width           float64       `json:"width"`
	Height          float64       `json:"height"`
	Angle           float64       `json:"angle"`
	StrokeColor     string        `json:"strokeColor"`
	BackgroundColor string        `json:"backgroundColor"`
	FillStyle       string        `json:"fillStyle"`
	StrokeWidth     int           `json:"strokeWidth"`
	StrokeStyle     string        `json:"strokeStyle"`
	Roughness       int           `json:"roughness"`
	Opacity         int           `json:"opacity"`
	GroupIds        []string      `json:"groupIds"`
	FrameId         *string       `json:"frameId"`
	Index           string        `json:"index"`
	Roundness       interface{}   `json:"roundness"`
	Seed            int64         `json:"seed"`
	Version         int64         `json:"version"`
	VersionNonce    int64         `json:"versionNonce"`
	IsDeleted       bool          `json:"isDeleted"`
	BoundElements   []interface{} `json:"boundElements"`
	Updated         int64         `json:"updated"`
	Link            *string       `json:"link"`
	Locked          bool          `json:"locked"`
	Text            string        `json:"text,omitempty"`
	FontSize        float64       `json:"fontSize,omitempty"`
	FontFamily      int           `json:"fontFamily,omitempty"`
	TextAlign       string        `json:"textAlign,omitempty"`
	VerticalAlign   string        `json:"verticalAlign,omitempty"`
	ContainerId     *string       `json:"containerId"`
	OriginalText    string        `json:"originalText,omitempty"`
	AutoResize      bool          `json:"autoResize,omitempty"`
	LineHeight      float64       `json:"lineHeight,omitempty"`
	Points          [][]float64   `json:"points,omitempty"`
	StartBinding    interface{}   `json:"startBinding"`
	EndBinding      interface{}   `json:"endBinding"`
	StartArrowhead  *string       `json:"startArrowhead"`
	EndArrowhead    *string       `json:"endArrowhead"`
	Polygon         bool          `json:"polygon"`
}

type AppState struct {
	GridSize              int                    `json:"gridSize"`
	GridStep              int                    `json:"gridStep"`
	GridModeEnabled       bool                   `json:"gridModeEnabled"`
	ViewBackgroundColor   string                 `json:"viewBackgroundColor"`
	LockedMultiSelections map[string]interface{} `json:"lockedMultiSelections"`
}

func Test(c *gin.Context) {

	var data ExcalidrawData

	shareableID := strings.Split("https://excalidraw.com/#json=4yfVrTy-LjzKWXlYv_tzx,Q3KWdNpRoSRGHS-TW0FL-Q", "json=")[1]

	decrypter := excalidrawdecrypt.CreateShareableExcalidrawDecrypter()

	plaintext, err := decrypter.Decrypt(shareableID)

	if err != nil {
		c.Status(500)
		return
	}

	err = json.Unmarshal([]byte(plaintext), &data)

	if err != nil {
		c.Status(500)
		return
	}

	c.JSON(200, data)

}
