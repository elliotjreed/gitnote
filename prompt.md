I keep all my notes in markdown files in a GitHub repository. I would like to have these organised and named according to the convention `yyyy-mm-dd note title.md`.

Act as a Go developer with experience in command line applications and unit testing. Create a cross-platform CLI application called "gitnote" which will manage the creation of a new note by executing `gitnote new`.

The `gitnote new` command should ask for a note category by scanning directories and listing them to be selected, if one directory is selected it should list the subdirectories. In all cases it should always have an option to create a new category. Creating a new category will create a new directory or subdirectory based on the context of the currently selected directory. It should then ask for a note title. The resulting new markdown file should be created with the current date in the format `yyyy-mm-dd note title.md` with the content being only the note title as a heading (`# note title`). Once create it should output the newly created file's relative path.

It should also have a command `gitnote index` to scan the directory and subdirectories for markdown notes and add them to readme.md in the form of a table of contents for each of the notes. The directories and subdirectories the notes are in should be headings in this file. For example, a file in `work/management/2025-01-05 managing expectations.md` should appear in the table of contents like:

```
## work

### management

[2025-01-05 managing expectations](/work/management/2025-01-05 managing expectations.md)
```

The `gitnote index` should validate the current readme.md against the directories contents and update the readme.md if necessary.

The command `gitnote search` should return a list of notes matching the title. For example `gitnote search "managing"` would return `work/management/2025-01-05 managing expectations.md`.

Adding the flag `gitnote search --full` should search for files containing the string in the file content as well as the title.

The command `gitnote commit` should add any new or updated files via `git`, with the commit message listing the new or changed files.

The command `gitnote pull` should pull in any changes from the remote git repository. If there are merge conflicts a message should be displayed with the option to manually fix the merge conflicts, or roll back. If there are no conflicts the new or updated files should be displayed.

Full test coverage is required. Clean code architecture should be followed.

