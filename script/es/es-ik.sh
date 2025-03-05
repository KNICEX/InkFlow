#!/bin/bash

# 1. 启动 Elasticsearch（后台进程）
/usr/local/bin/docker-entrypoint.sh &

# 2. 等待 ES 启动完成
echo "Waiting for Elasticsearch to be ready..."
until curl -s -u "elastic:changeme" http://localhost:9200 | grep -q "You Know, for Search"; do
  sleep 5
done

if [ -d "plugins/analysis-ik" ]; then
  echo "IK plugin is already installed. Skipping installation."
else
  echo "Elasticsearch is ready. Installing IK plugin..."
  bin/elasticsearch-plugin install --batch https://get.infini.cloud/elasticsearch/analysis-ik/8.17.2

  echo "IK plugin installed. Restarting Elasticsearch..."

fi

# 3. 终止后台 Elasticsearch 进程
pkill -f elasticsearch

# 4. 重新启动 Elasticsearch（前台运行，保持容器存活）
exec /usr/local/bin/docker-entrypoint.sh