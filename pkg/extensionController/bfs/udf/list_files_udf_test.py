#pylint: disable=missing-function-docstring,missing-module-docstring,missing-class-docstring
from pathlib import Path
import os
import pytest
import list_files_udf

#pylint: disable=too-few-public-methods # Mock does not need more methods
class ExaContextMock:
    path: str
    emitted_rows: list
    def __init__(self, path:str) -> None:
        self.path = path
        self.emitted_rows = []

    def emit(self, *args) -> None:
        self.emitted_rows.append(args)

def run(context: ExaContextMock) -> None:
    list_files_udf.run(context)

def run_get_emitted_rows(bfs_path: Path) -> list:
    context = ExaContextMock(str(bfs_path))
    list_files_udf.run(context)
    return context.emitted_rows

def create_file(path: Path, content: str) -> None:
    os.makedirs(path.parent, exist_ok=True)
    with open(path, mode="w", encoding="UTF-8") as f:
        f.write(content)

def test_fails_for_empty_path() -> None:
    with pytest.raises(ValueError) as exception:
        run(ExaContextMock(""))
    assert "Argument 'path' not defined" in str(exception)

def test_empty_dir(tmp_path: Path) -> None:
    assert len(run_get_emitted_rows(tmp_path)) == 0

def test_single_file(tmp_path: Path) -> None:
    file1 = tmp_path/"file1.txt"
    create_file(file1, "content")
    rows = run_get_emitted_rows(tmp_path)
    assert rows == [("file1.txt", str(file1), 7)]

def test_multiple_files(tmp_path: Path) -> None:
    file1 = tmp_path/"file1.txt"
    file2 = tmp_path/"file2.txt"
    create_file(file1, "content")
    create_file(file2, "even more content")
    rows = run_get_emitted_rows(tmp_path)
    assert (len(rows) == 2
            and ("file1.txt", str(file1), 7) in rows
            and ("file2.txt", str(file2), 17) in rows)

def test_sub_dir(tmp_path: Path) -> None:
    file1 = tmp_path/"dir1"/"file1.txt"
    file2 = tmp_path/"dir2"/"file2.txt"
    create_file(file1, "content")
    create_file(file2, "even more content")
    rows = run_get_emitted_rows(tmp_path)
    assert (len(rows) == 2
            and ("file1.txt", str(file1), 7) in rows
            and ("file2.txt", str(file2), 17) in rows)

def test_skips_dir_with_exaudf_subdir(tmp_path: Path) -> None:
    ignored_dir = tmp_path/"ignored"
    file1 = ignored_dir/"ignored.txt"
    file2 = tmp_path/"file.txt"
    create_file(file1, "content")
    create_file(file2, "even more content")
    os.mkdir(ignored_dir/"exaudf")

    rows = run_get_emitted_rows(tmp_path)
    assert rows == [("file.txt", str(file2), 17)]
