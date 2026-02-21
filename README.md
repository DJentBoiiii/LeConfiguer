`le add {path-to-file-file} -t {type} -n {name} -i {id} -e {environment}` - upload config to MinIO
`le remove -i {id}` - delete config from MinIO and indexing
`le update -i {id} {file}` - update config
`le diff -i {id} {file}` - diff with previous version if exists
`le use -i {id}` - download file. Also example use: `cat <(le use -i {id})`. For custom file name `le use -i {id} -o "filename"`
`le rollback -i {id}` - rollback config to previous version