local Discordia = require'discordia'
local Client = Discordia.Client()
local Http = require'coro-http'
local sfind,slower,ssub = string.find,string.lower,string.sub
local request = Http.request 
function writefile(File,...)
    local File = io.open(File,"w")
    File:write(...)
    if File then -- error handling in edge-cases
        File:close()
    end
end

function readfile(File)
    local File = io.open(File,"r")
    if File then 
        local Contents = File:read()
        File:close()
        return Contents
    end 
end 
if not readfile("Token.txt") then writefile("Token.txt") end

local function testnitro(Code)
    local Head,Body = request("GET","https://ptb.discordapp.com/api/v6/entitlements/gift-codes/"..Code,{{"Authorization",Token}})
    return Head.code == 200 and true
end

local function redeemcode(Code)
    local Head,Body = request("POST","https://ptb.discordapp.com/api/v6/entitlements/gift-codes/"..Code.."/redeem",{{"content-type","application/json"},{"Authorization",Token}},'{"payment_source_id":null}')
    return Head.code
end

Client:on('messageCreate',function(Message)
local Content = Message.content
local Find = sfind(slower(Content),"discord.gift/")
    if Find then
    local Gift = ssub(Content,Find + 13)
        if #Gift == 16 and testnitro(ssub(Gift,1,16)) then 
            if redeemcode(ssub(Gift,1,16)) == 200 then  -- a simple check
                print('Claimed '..ssub(Gift,1,24)..' from '..Message.author.tag)
           end
        elseif #Gift == 24 and testnitro(ssub(Gift,1,24)) then 
           if redeemcode(ssub(Gift,1,24)) == 200 then  -- a simple check
                print('Claimed '..ssub(Gift,1,24)..' from '..Message.author.tag)
           end
        end
    end
end)

if not readfile("Token.txt") and not args[2] then error'You need to put your token as a second argument or put it in the Token.txt file.' end
Client:run(readfile("Token.txt") or args[2])
