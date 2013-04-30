/* Copyright 2013 Robert Zaremba
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package surfer

import (
	"fmt"
	"github.com/scale-it/go-log"
)

type Error struct {
	Level log.Level
	E     error
	Msg   string
}

func (e Error) Error() string {
	return fmt.Sprintf("[%d] %s; %s", e.Level, e.Msg, e.E)
}

func (e Error) Log(l *log.Logger) {
	l.Logf(e.Level, "%s; %s", e.Msg, e.E)

}
