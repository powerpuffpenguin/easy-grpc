Target="server"
Docker="server"
Dir=$(cd "$(dirname $BASH_SOURCE)/.." && pwd)
Version="v0.0.1"
View=0
Platforms=(
    darwin/amd64
    windows/amd64
    linux/arm
    linux/amd64
)
UUID="bb888120-b359-11ed-b2eb-bdc03e56513e"
Protos=(
    system/system.proto
    session/session.proto
    user/user.proto
    logger/logger.proto
)