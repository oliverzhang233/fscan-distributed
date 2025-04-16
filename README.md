# Fscan 分布式扫描系统
项目概述
本项目是一个基于 Fscan 的分布式扫描系统，借助在子公司内网部署 Agent 探针来下发扫描任务，扫描完成后收集结果并上报至总部服务器，总部服务器可同时查看所有扫描结果。系统采用分布式架构，具备任务下发、扫描执行、结果收集和结果展示等功能，满足安全性、可靠性和可扩展性的需求。

系统架构
总部服务器：负责创建、下发和管理扫描任务，接收并存储各个子公司的扫描结果，同时提供可视化界面供管理员查看。
子公司 Agent 探针：接收总部服务器下发的扫描任务，使用 Fscan 工具执行扫描，并将结果上报给总部服务器。
技术选型
总部服务器：
操作系统：Linux（如 CentOS、Ubuntu）
编程语言：Go
Web 框架：标准库 net/http
数据库：暂未使用，可根据需求扩展
消息队列：RabbitMQ
子公司 Agent 探针：
操作系统：根据子公司内网环境选择，如 Windows 或 Linux
编程语言：Go
Fscan 工具：开源的 Fscan 项目

目录结构
plaintext
fscan-system/
├── server/
│   ├── main.go
│   ├── server_config.yaml
│   └── deploy_server.sh
├── agent/
│   ├── main.go
│   ├── agent_config.yaml
│   └── deploy_agent.sh
└── README.md
安装与部署
服务端部署
克隆项目
bash
git clone https://github.com/your-repo/fscan-system.git
cd fscan-system/server

配置服务端
编辑 server_config.yaml 文件，根据实际情况修改消息队列和服务器端口配置：
yaml
amqp_url: "amqp://guest:guest@localhost:5672/"
task_queue: "scan_tasks"
result_queue: "scan_results"
server_port: 8080
一键部署
运行部署脚本：
bash
./deploy_server.sh
Agent 端部署
克隆项目
bash
git clone https://github.com/your-repo/fscan-system.git
cd fscan-system/agent

配置 Agent 端
编辑 agent_config.yaml 文件，根据实际情况修改消息队列配置：
yaml
amqp_url: "amqp://guest:guest@localhost:5672/"
task_queue: "scan_tasks"
result_queue: "scan_results"
放置 Fscan 工具
将 fscan 可执行文件放置在与 agent 可执行文件相同的目录下，或者放置在系统 PATH 环境变量包含的目录中，如 /usr/local/bin。

一键部署
运行部署脚本：
bash
./deploy_agent.sh

使用方法
任务下发
向服务端的 /send_task 接口发送 POST 请求，请求体为 JSON 格式，包含扫描任务信息，例如：
json
{
    "host": "192.168.1.1"
}

结果查看
服务端接收到 Agent 端上报的扫描结果后，会将结果打印到日志中。你可以根据需求扩展服务端代码，将结果存储到数据库并提供可视化界面进行查看。

注意事项
确保消息队列（RabbitMQ）正常运行，并且服务端和 Agent 端能够连接到消息队列。
在 Agent 端部署时，确保 fscan 工具可正常执行，并且具有相应的权限。
为了保证系统的安全性，建议在生产环境中对消息队列和服务端进行安全配置，如设置用户名、密码和访问控制。

贡献与反馈
如果你在使用过程中遇到问题或有任何建议，欢迎在 GitHub 上提交 Issue 或 Pull Request。

许可证
本项目采用 MIT 许可证。
