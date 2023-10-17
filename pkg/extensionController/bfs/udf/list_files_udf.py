#pylint: disable=missing-function-docstring,missing-module-docstring,missing-class-docstring

from pathlib import Path
import stat
import os
from dataclasses import dataclass
from typing import Generator

@dataclass(frozen=True)
class ResultPath:
    path: Path
    _stat: os.stat_result

    def get_name(self) -> str:
        return self.path.name
    def get_absolute_path(self) -> str:
        return str(self.path)
    def get_size(self) -> int:
        return self._stat.st_size
    def is_file(self) -> bool:
        return stat.S_ISREG(self._stat.st_mode)
    def is_dir(self) -> bool:
        return stat.S_ISDIR(self._stat.st_mode)


def run(ctx) -> None:
    if not ctx.path:
        raise ValueError("Argument 'path' not defined")
    for f in list_recursively(Path(ctx.path)):
        ctx.emit(f.get_name(), f.get_absolute_path(), f.get_size())

def accept_file(file: ResultPath) -> bool:
    if file.is_dir():
        udf_dir = file.path/"exaudf"
        return not os.path.isdir(udf_dir)
    return file.is_file()

def list_dir(path: Path) -> Generator[ResultPath, None, None]:
    try:
        for f in path.iterdir():
            try:
                file = ResultPath(f, os.lstat(f))
                if accept_file(file):
                    yield file
            except FileNotFoundError:
                pass
    except PermissionError:
        pass

def list_recursively(path: Path) -> Generator[ResultPath, None, None]:
    for child in list_dir(path):
        if child.is_file():
            yield child
        if child.is_dir():
            yield from list_recursively(child.path)
