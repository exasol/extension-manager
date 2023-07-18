import pathlib
import stat

def run(ctx):
	files = [p for p in pathlib.Path(ctx.path).rglob("*") if p.exists() and "EXAClusterOS" not in p.parts]
	files = [read_path(p) for p in files]
	files = [p for p in files if p["is_file"]]
	for file in files:
		ctx.emit(file['name'], file['path'], file['size'])

def read_path(path):
	file_stat = path.stat()
	size = file_stat.st_size
	is_file = stat.S_ISREG(file_stat.st_mode)
	return {'name': path.name, 'path':str(path), 'size':size, 'is_file':is_file}
