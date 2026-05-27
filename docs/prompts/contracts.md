查看 contracts目录，进行日志更新或者对shema进行修改。

之后将内容进行git add后，给出类似于如下形式的：
```shell
git commit -m "$(cat <<'EOF'
<commit message>

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```
