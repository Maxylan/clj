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

## Examples

**View** ticket `PROJ-1337`
<br/>`clj PROJ-1337`

**Find Transitions** for tickets `PROJ-1337` and `CUST-128`
<br/>`clj PROJ-1337 CUST-128`

**Find Users** matching search "`Maxy`"
<br/>`clj users Maxy`

**Assign** user "Maxylan" to ticket `PROJ-1337`, leave a **comment** on `PROJ-1337` _(print the comment section)_, and finally **print** a detailed **view** of ticket `PROJ-1337` _(all in a single command)_
<br/>`clj set assignee "Maxylan" on PROJ-1337 and comment "Assigning ticket to myself" on PROJ-1337 -o and view PROJ-1337 --detailed`

## Help

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
├ clj [init|setup]
╰ "Performs initial setup / configuration. Can be re-run later to reconfigure this utility."

» Search / Find matching users
├ clj [u|user|users] "<Value>"
├ "Search for users matching <Value>, remember to use quotes if your search contains whitespace characters."
╰ Subcommands:
        • Detailed view
        ├ clj <...> [-d|--detailed]
        ╰ "Include as many details as possible"
        • Exact / Strict matching
        ├ clj <...> [-e|--exact]
        ╰ "Force exact / strict matching."
        • Include Deleted Users
        ├ clj <...> [-i|--include-deleted]
        ╰ "Override default behaviour that filters out deleted users."


» Available Ticket Transitions (statuses)
├ clj [stat|statuses|transitions] on <Tickets...>
├ "Retrieves information about each available transition (status), for each given ticket."
╰ Subcommands:
        • Detailed view
        ├ clj <...> [-d|--detailed]
        ╰ "Include as many details as possible"


» Set some field(s) on Ticket(s)
├ clj set <Field> <Value> on <Tickets...>
├ "Update some field on ticket(s), see subcommands. Ex. `clj set status "Done" on PROJ-1337`"
╰ Subcommands:
        • Field: assignee
        ├ clj set assignee <User> on <...>
        ├ "Set 'Assignee' of ticket(s). Picks best user-match. See subcommand `users` for available users. Ex. `clj users`"
        ╰ Subcommands:
                • Exact / Strict matching
                ├ clj <...> [-e|--exact]
                ╰ "Force exact / strict matching in your user search."
                • Prompt to select
                ├ clj <...> [-s|--select]
                ╰ "Instead of auto-picking best match, this will prompt you for input. Let's you pick the user."

        • Field: reporter
        ├ clj set reporter <User> on <...>
        ├ "Set 'Reporter' of ticket(s). See subcommand `users` for available users. Ex. `clj users`"
        ╰ Subcommands:
                • Exact / Strict matching
                ├ clj <...> [-e|--exact]
                ╰ "Force exact / strict matching in your user search."
                • Prompt to select
                ├ clj <...> [-s|--select]
                ╰ "Instead of auto-picking best match, this will prompt you for input. Let's you pick the user."

        • Field: status
        ├ clj set [status|transition] <Value> on <...>
        ╰ "Set 'status' ('transition') on ticket(s). See subcommand `stat` for available transitions. Ex. `clj stat <PROJ-1337>`"
        • Print Updated Tickets
        ├ clj set <...> on <...> [-p|--print]
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
