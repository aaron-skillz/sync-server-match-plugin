--[[
 Copyright 2019 The Nakama Authors

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
--]]

local nk = require("nakama")


local function create_match(context, payload)
  local decoded = nk.json_decode(payload)
  local params = {
    skillz_match_id = decoded.match_id,
    debug = true,
    label = "skillz-match-sync"
  }

  local module = "match"
  local match_id = nk.match_create(module, params)
  return nk.json_encode({ match_id = match_id })
end

nk.register_rpc(create_match, "create_skillz_match")
