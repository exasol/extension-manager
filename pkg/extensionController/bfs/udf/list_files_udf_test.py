#pylint: disable=missing-function-docstring,missing-module-docstring,missing-class-docstring
import sys
from pathlib import Path
import pytest
import list_files_udf
from typing import Any


class ExaContextMock:
    path: str
    emitted_rows: list
    def __init__(self, path:Path) -> None:
        self.path = str(path)
        self.emitted_rows = []

    def emit(self, *args) -> None:
        print("emit:", args)
        self.emitted_rows.append(args)

def run(context: ExaContextMock) -> None:
    list_files_udf.run(context)

def run_get_emitted_rows(bfs_path: Path) -> list:
    context = ExaContextMock(bfs_path)
    list_files_udf.run(context)
    return context.emitted_rows


def create_file(path: Path, content: str) -> None:
    with open(path, mode="w", encoding="UTF-8") as f:
        f.write(content)


def test_run_empty_dir(tmp_path: Path) -> None:
    assert len(run_get_emitted_rows(tmp_path)) == 0

def test_run_single_file(tmp_path: Path) -> None:
    file1 = tmp_path/"file1.txt"
    create_file(file1, "content")
    rows = run_get_emitted_rows(tmp_path)
    assert len(rows) == 1 and rows[0] == ("file1.txt", str(file1), 7)
