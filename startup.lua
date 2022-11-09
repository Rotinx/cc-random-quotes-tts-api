local dfpwm = require("cc.audio.dfpwm")
local speaker = peripheral.find("speaker")


local decoder = dfpwm.make_decoder()

local quote = http.get("http://ip/get", {}, true)
local res = quote.readAll()

local f = assert(io.open('latest.dfpwm', 'wb')) -- open in "binary" mode
f:write(res)
f:close()

for chunk in io.lines("latest.dfpwm", 16 * 1024) do
    local buffer = decoder(chunk)

    while not speaker.playAudio(buffer) do
        os.pullEvent("speaker_audio_empty")
    end
end