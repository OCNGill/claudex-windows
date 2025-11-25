In Progress:
  - 

To Do: 
  - Adjust Agents' models:
    - Documenter hook: haiku
    - Architect: opus
    - Researcher: haiku/sonnet
    - Engineer: sonnet/opus
    - Prompt Engineer: sonnet/opus
  - Analyze all hooks via loggers: /hooks (look into Stop, SubagentStart)
    - Documenter:
      - Analyze if documenter could be better if triggered on PreCompact or Stop
      - Idea: With autocompact disabled, when context is full you get an error. Use gemini to read the transcript and produce or update the overview document. Remove from the transcript tools executions and similar non-relevant things (like I did with claudex already). It's important that gemini produce several docs, being the overview a small one with pointers to be loaded only on demand.
  - Add engineer with flexible skills (typescript, python, etc)
  - Improve researcher profile to generate more concise documentation
  - Improve architect profile to generate more concise documentation
  - Find a solution to allow the user to define where the relevant documentation is located (product, architecture, standards, etc)
  - Review all agents to adjust their output format to make sure they provide enough context to the caller (team lead) but avoid verbose responses