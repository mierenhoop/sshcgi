#!/usr/bin/lua

-- TODO: allow script to accept/reject connections
-- by checking username/password/key

-- Go to new line is done with \r\n in raw terminal
io.write("Started process\r\n")
io.write("PATH="..os.getenv"PATH".."\r\n")
io.flush() -- Flushing has to be done manually

while true do
    local c = io.read(1)
    -- If C-c or C-d exit
    if c == "\3" or c == "\4" then break end
    io.write(string.format("You wrote %q.\r\n", c))
    io.flush()
end

io.write("Exiting...\r\n")
io.flush()