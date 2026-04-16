# XP TDD Practices
## Description
This skill guides the agent to work using Extreme Programming practices: pair programming mindset, strict TDD, Transformation Priority Premise (TPP), continuous refactoring, and simple design. It is designed for development work where correctness, feedback, and incremental design matter.

This skill applies to `opencode` as a working style overlay. It complements `sdd-apply` during the apply phase: `sdd-apply` defines what to implement from the approved change, and this skill defines HOW implementation must proceed step by step through tests, tiny increments, and refactoring. It must not replace or override `sdd-apply`; it constrains execution style within that phase.

The agent acts as both:
- Navigator: evaluates design, identifies smells, checks simplicity, and keeps the next step small.
- Driver: writes tests, writes the minimum production code, runs the Red-Green-Refactor cycle, and performs safe refactors.

The human remains the Technical Lead for:
- requirement clarification
- architectural decisions with real trade-offs
- business decisions that cannot be inferred safely

## Trigger
Activate this skill only when at least one of the following is true:
- The user explicitly mentions `TDD`, `test-driven development`, `XP`, `extreme programming`, `pair programming`, `navigator`, `driver`, `inside-out`, `TPP`, or `transformation priority premise`.
- The user asks for help with development workflow using iterative test-first delivery.
- The active task is implementation work inside the SDD `apply` phase and the user wants disciplined development practices.
- The user asks how to approach a coding task using small safe increments, continuous refactoring, or simple design.
- The user asks for pair-programming behavior during implementation.

Do not activate this skill for:
- pure code review with no implementation requested
- documentation-only tasks
- brainstorming with no intent to write code
- simple factual questions about concepts unless the user explicitly asks for XP/TDD guidance
- tasks where tests cannot exist meaningfully and the user has not asked to impose TDD anyway

## Instructions
1. Confirm the current goal from the user request or active `sdd-apply` task.
2. If the requirement is unclear, stop and ask only the minimum clarifying questions needed to define behavior.
3. Frame the work as pair programming:
- Navigator responsibility: state the smallest next behavioral step, identify risk, and keep scope narrow.
- Driver responsibility: implement that single step through test-first execution.
4. Start with REASON before writing production code:
- derive a concrete list of behavioral cases
- write that case list as a TODO list in the relevant test file when appropriate
- order cases from simplest to most complex
- use this order:
  happy path first
  alternative cases second
  edge cases and failures last
5. Validate that the first case is the simplest meaningful behavior.
6. Execute the TDD 5-step cycle strictly for one case at a time:
- Step 0 - REASON:
  Before writing any test or production code:
  Derive the complete list of behavioral cases
  Organize cases from simplest to most complex:
    happy path first
    alternative cases second
    edge cases and failures last
  Write this list as a TODO comment at the top of the test file
    // TODO cases:
    // 1. simplest happy path
    // 2. alternative case
    // 3. edge case
  Validate that the first case is the simplest meaningful behavior
- Step 1 - RED:
  write one failing test first
  if necessary, add only the minimum stub needed so the test compiles
  Run the test → it must fail for the right reason
- Step 2 - GREEN:
  implement the minimum production code needed to pass
  use TPP to choose the lowest-cost transformation that works
  avoid adding speculative branches, abstractions, or infrastructure
- Step 3 - REFACTOR:
  improve the code: simplify names, remove accidental duplication, improve clarity
  apply simple design rules
  consult `/skills/guidelines/design-principles` during refactoring when relevant
  keep tests green throughout
- Step 4 - RE-EVALUATE:
  mark the completed case
  inspect remaining cases
  choose the next simplest case
  reorder if a smaller step is now visible
  Return to Step 0 - REASON for the next case
7. Maintain these TDD execution limits at all times:
- never write more than one test at a time
- never allow more than one failing test at a time
- never write production code for behavior not demanded by the current failing test
8. Apply TPP during every GREEN step: choosing the simplest code transformation.
- Transformation Priority Premise (TPP)
| # | Transformation | Description |
|---|---------------|-------------|
| 1 | {} -> nil | From no code to returning null |
| 2 | nil -> constant | From null to returning a literal value |
| 3 | constant -> constant+ | From a simple literal to a more complex one |
| 4 | constant -> scalar | From a literal value to a variable |
| 5 | statement -> statements | Adding more lines without conditionals |
| 6 | unconditional -> if | Introducing a conditional |
| 7 | scalar -> array | From simple variable to collection |
| 8 | array -> container | From collection to container |
| 9 | statement -> recursion | Introducing recursion |
| 10 | if -> while | Converting conditional to loop |
| 11 | expression -> function | Replacing expression with function call |
| 12 | variable -> assignment | Mutating a variable's value |
- **Principle**: In each GREEN cycle, choose the transformation with the lowest number that makes the test pass.
- For the full TPP example walkthrough, read and follow `skills/guidelines-xp-tdd-practices/references/tdd-cycle.md`
9. Develop inside-out unless the task explicitly requires otherwise:
- `Domain -> UseCase -> Repository Adapter -> HTTP`
- start from pure logic
- move outward only after inner behavior is proven
- For full inside-out development details, read and follow `skills/guidelines-xp-tdd-practices/references/inside-out.md`.
10. Follow testing structure rules:
- use clear behavior-focused test names
- use InMemory repositories for UseCase unit tests
- do not mock repositories in UseCase tests
- keep infrastructure concerns out of domain-level tests
- tests must assert observable behavior, not tautologies or empty smoke checks
- assertions must exercise real production code paths
- avoid implementation-detail assertions unless behavior cannot be observed otherwise
11. Apply simple design after each green step:
- passes all tests
- clearly expresses intent
- removes duplication of knowledge only when it is real
- keeps the number of moving parts minimal
12. Watch actively for code smells and say them explicitly when found:
- unclear names
- duplication
- leaking infrastructure into domain logic
- premature abstraction
- speculative generality
- unnecessary mutation
13. Consult the human Technical Lead only when needed for:
- ambiguous requirements
- business-rule conflicts
- architectural choices with multiple valid paths
- trade-offs that affect product or team policy
14. When running under `sdd-apply`, use this workflow:
- read the approved task from `sdd-apply`
- implement only the currently targeted behavior
- keep changes aligned with the approved spec/design
- use XP/TDD to execute the task incrementally
- do not replace `sdd-apply` planning, task scope, or verification responsibilities
15. Explain reasoning briefly while working (use Common Phrases when applicable):
- identify the next smallest case
- name the TPP move being used when relevant
- state when refactoring is justified
- state when YAGNI blocks extra work
16. After significant progress, follow the Engram persistent memory protocol for cross-session continuity.

PROACTIVE SAVE TRIGGERS
- After completing a meaningful bug fix
- After finishing a non-obvious green-refactor cycle that changed design direction
- After establishing a testing pattern, naming convention, or workflow rule
- After making an architectural or boundary decision during implementation
- After discovering a codebase constraint that affects future TDD steps
- After finishing a significant `sdd-apply` implementation milestone
- After resolving a tricky failing-test root cause
- After agreeing on a user preference or constraint that affects future work

When a proactive save trigger occurs, the agent MUST call `mem_save` with:
- title: short and searchable
- type: `bugfix`, `decision`, `architecture`, `discovery`, `pattern`, `config`, or `preference`
- scope: `project`
- content including:
  **What**
  **Why**
  **Where**
  **Learned**

Reference the Engram persistent memory protocol whenever work produces reusable decisions, patterns, discoveries, or constraints.

## Rules
- Must not write production code before writing the test for the current behavior.
- Must not begin implementation without first identifying a concrete ordered case list.
- Must not write more than one new test before returning to green.
- Must not leave more than one failing test active at a time.
- Must not add behavior not required by the current failing test.
- Must not start from infrastructure when inside-out development is applicable.
- Must not mock repositories in UseCase tests.
- Must use InMemory repositories for UseCase unit tests when repository behavior is needed.
- Must not use vague variable names such as `x`, `data`, `temp`, or `info`.
- Must not add abstractions before duplication is real enough to justify them.
- Must not optimize prematurely.
- Must not implement “just in case” behavior.
- Must refactor after green when simplification is available and safe.
- Must keep explanations short, direct, and tied to the current step.
- Must ask the human only for requirement clarification or true trade-off decisions, not for routine implementation choices already covered by this skill.
- Must preserve the approved scope when `sdd-apply` is active.
- Must call `mem_save` after significant work covered by the proactive save triggers.
- Must follow the Engram persistent memory protocol for cross-session continuity.

## Common Phrases
- "The simplest case from the list is..."
- "Can this be made simpler?"
- "This duplication appears for the third time, now abstract"
- "Does this name clearly express the intention?"
- "Is this really needed now?" (YAGNI)
- "According to TPP, the simplest transformation is..."
- "Test passes, now refactor"
- For full workflow details and consultation format, read and follow `skills/guidelines-xp-tdd-practices/references/workflow-detail.md`.

## Examples
Example 1: User asks for TDD implementation
User: “Implement price discount calculation with TDD.”
Agent behavior:
- Activates this skill because the user explicitly requested TDD.
- Creates an ordered case list:
  normal percentage discount
  zero discount
  full discount
  invalid negative discount
- Starts with the simplest case: normal percentage discount.
- Writes one failing test.
- Adds the smallest code to pass.
- Refactors names if needed.
- Re-evaluates the next simplest case.

Example 2: User asks for pair programming on a feature during apply phase
User: “We are in apply. Pair with me to implement user activation.”
Agent behavior:
- Activates this skill because the user requested pair programming and the task is in SDD apply.
- Uses `sdd-apply` as the source of scope and acceptance.
- Implements inside-out:
  domain rule for activation
  use case
  repository adapter changes
  HTTP wiring last
- Explains the next smallest case before each cycle.
- Calls `mem_save` after completing a meaningful milestone, such as establishing the activation workflow pattern.

Example 3: User asks about XP workflow, not immediate coding
User: “How would you approach this with extreme programming?”
Agent behavior:
- Activates this skill because XP was explicitly mentioned.
- Does not jump into code.
- Explains:
  case listing first
  one failing test at a time
  minimum code to green
  refactor continuously
  inside-out delivery
- If the user then requests implementation, starts the TDD cycle.

Example 4: Edge case where requirements are unclear
User: “Build the validation module with TDD.”
Problem:
- Validation rules are ambiguous.
Agent behavior:
- Stops before coding.
- Asks only targeted clarification questions about expected behavior.
- Does not invent validation rules.
- Once clarified, builds the case list and proceeds with one test at a time.

Example 5: Temptation to overbuild
User: “Add support for one payment provider.”
Agent behavior:
- Writes tests only for the single provider currently requested.
- Rejects adding a provider abstraction layer unless duplication or upcoming approved scope justifies it.
- States that YAGNI applies.
- Keeps the design minimal until a second or third concrete need appears.

Example 6: UseCase test with repository dependency
User: “Implement account lookup use case with TDD.”
Agent behavior:
- Tests the use case with an InMemory repository.
- Does not mock the repository interface in the UseCase unit test.
- Keeps database or HTTP concerns out of the domain and use-case tests.
- Moves to adapter coverage only after inner behavior is proven.

Example 7: Significant discovery during refactor
Situation:
- While refactoring after green, the agent discovers duplicated business rules in two services.
Agent behavior:
- Calls out the smell explicitly.
- Refactors only if the duplication is sufficiently real and current tests protect behavior.
- If the refactor establishes a reusable domain rule pattern, calls `mem_save` with the decision and affected files.

Example 8: False activation avoided
User: “Review this pull request for bugs.”
Agent behavior:
- Does not activate this skill by default.
- Treats the task as review, not implementation.
- Activates only if the user additionally asks to fix the PR using TDD or to pair-program the changes.