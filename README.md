# Command-line Jira (clj)

Can be used to interact with Jira via the command-line.

I love creating tools for myself as a way to practice a new technology!

## Features

 _Note: Most, if not all commands, can be passed more than one ticket at a time! Just separate them w/ white-space_

 * Viewing ticket(s) and their comments
 * Chaining multiple commands w/ `.. and ..`
 * Parsing Jira markup to a pretty console print
 * Converting some simple markup to Jira markup
 * Post comment on ticket(s)
 * Get transitions _(statuses)_ of ticket(s)
 * "Transition" the status of ticket(s)

## Usage

**TLDR;** `clj` + any number of ticket IDs you want to view _(separated by spaces)_.

Configuration happens automagically, as you try to use the tool for the first time. If you ever need to re-configure, run `clj init`

For more details about capabilities, `clj help` will give you this lovely print, explaining each subcommand in detail.

```
$ clj
» Help
├ clj help | alt. [-h|--help] after any command.
├ "Prints usage & detailed explanations of each subcommand of their accepted arguments"
╰ Subcommands:
        • Minimal
        ├ clj help minimal|[-m|--minimal]
        ╰ "Omit lengthy descriptions, only printing command name + usage."


» Init (initial setup / configuration)
├ clj init
╰ "Performs initial setup / configuration. Can be re-run later to reconfigure this utility."

» Set some field(s) on Ticket(s)
├ clj set <Field> <Value> on <Tickets...>
├ "Update some field on ticket(s), see subcommands. Ex. `clj set status "Done" on PROJ-1337`"
╰ Subcommands:
        • Field: status
        ├ clj set status <Value> on <...>
        ╰ "Set 'status' field on ticket(s). Statuses defined in a project's workflow, see `clj statuses <PROJ-1337>`"
        • View Updated Tickets
        ├ clj set <...> on <...> [-t|--view-tickets]
        ╰ "Print updated ticket(s)."


» Comment on Ticket(s)
├ clj comment "Lorem ipsum dolor.." on <Tickets...>
├ "Creates a new comment on each given ticket. Outputs their comment sections."
╰ Subcommands:
        • Oldest First
        ├ clj <...> [--oldest-first]
        ╰ "Change default sort-order to show the oldest comments first."


» View Ticket(s)
├ clj <Tickets...>
├ "Retrieves each given ticket, printing formats and prints each one to the console."
╰ Subcommands:
        • Detailed view
        ├ clj <Tickets...> [-d|--detailed]
        ╰ "Include as many details as possible"
        • Include comments
        ├ clj <Tickets...> [-c|--comments]
        ╰ "Render the whole comment section"
        • Only comments
        ├ clj <Tickets...> [-o|--only-comments]
        ╰ "Render *only* the comment section"
        • Oldest First
        ├ clj <Tickets...> [--oldest-first]
        ╰ "Change default sort-order to show the oldest comments first. Assumes [-c|-o]"
```
