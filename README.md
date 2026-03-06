# Command-line Jira (clj)

Can be used to interact with Jira via the command-line.

I love creating tools for myself as a way to practice a new technology!

## Usage

**TLDR;** `clj` + any number of ticket IDs you want to view _(separated by spaces)_.

Configuration happens automagically, as you try to use the tool for the first time. If you ever need to re-configure, run `clj init`

For more details about capabilities, `clj help` will give you this lovely print, explaining each subcommand in detail.

```
clj v1.0.0 [clj]

» Help
├ clj help | alt. [-h|--help] after any command.
├ "Prints usage & detailed explanations of each subcommand of their accepted arguments"
╰ Subcommands:
        • Minimal
        ├ clj help minimal|[-m|--minimal]
        ╰ Omit lengthy descriptions, only printing command name + usage.


» Init (initial setup / configuration)
├ clj init
╰ "Performs initial setup / configuration. Can be re-run later to reconfigure this utility."

» Comment on Ticket(s)
├ clj [c|comment] <PROJ-1337> <...> "Lorem ipsum dolor.."
╰ "Creates a new comment on each given ticket, outputs their comment sections."

» View Ticket(s)
├ clj <PROJ-1337> <...>
├ "Retrieves each given ticket, printing their title + description."
╰ Subcommands:
        • Detailed view
        ├ clj <...> [-d|--detailed]
        ╰ Include as many details as possible
        • Include comments
        ├ clj <...> [-c|--comments]
        ╰ Render the whole comment section
        • Only comments
        ├ clj <...> [-o|--only-comments]
        ╰ Render *only* the comment section
```
