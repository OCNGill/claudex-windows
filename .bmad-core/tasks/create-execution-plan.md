<!-- Powered by BMAD™ Core -->

# Create Execution Plan Document

## ⚠️ CRITICAL EXECUTION NOTICE ⚠️

**THIS IS AN EXECUTABLE WORKFLOW - NOT REFERENCE MATERIAL**

When this task is invoked:

1. **MANDATORY CLARIFICATION PHASE FIRST** - You MUST start with an explicit clarification phase before ANY document creation
2. **CLARIFY BEFORE DOCUMENTING** - Resolve ALL questions with the user before producing document sections
3. **USE MCP TOOLS** - Query documentation via context7 MCP, use sequential-thinking for complex analysis
4. **FINAL DECISIONS ONLY** - Documents must contain ONLY final decisions, not alternatives or rationale discussions
5. **MANDATORY STEP-BY-STEP EXECUTION** - Each section must be processed sequentially with user feedback
6. **ELICITATION IS REQUIRED** - When `elicit: true`, you MUST use the 1-9 format and wait for user response

**VIOLATION INDICATOR:** If you create document sections with alternatives, options, or rationale before clarifying with user, you have violated this workflow.

**WORKFLOW VIOLATION:** Starting document creation before completing the explicit clarification phase violates this task.

## PHASE 1: MANDATORY CLARIFICATION PHASE

**THIS PHASE MUST BE COMPLETED BEFORE ANY DOCUMENT CREATION BEGINS**

### Overview

When the user requests an execution plan (for a new feature, refactoring, or any code work), you MUST start with this explicit clarification phase. The goal is to ensure complete understanding and alignment BEFORE creating any document sections.

### Step-by-Step Clarification Process

**Step 1: Acknowledge and Set Context**

- Acknowledge the user's request
- Briefly summarize what you understand they want to accomplish
- State explicitly: "Before I create the execution plan, I need to clarify some details to ensure we're aligned."

**Step 2: Gather Context**

- Read any provided PRD, architecture documents, user stories, or requirements
- Search the codebase for relevant existing implementations
- Identify all libraries, SDKs, frameworks, third-party services mentioned

**Step 3: Query Documentation** (MANDATORY)

- Use `mcp__context7__resolve-library-id` to identify libraries
- Use `mcp__context7__get-library-docs` to fetch up-to-date documentation
- Examples: Firebase Functions, OpenAI SDK, TypeScript, React, etc.

**Step 4: Analyze Complexity** (when applicable)

- Use `mcp__sequential-thinking__sequentialthinking` for:
  - Complex architectural decisions
  - Multiple alternative approaches
  - Interconnected system changes
  - Trade-off analysis

**Step 5: Ask Clarifying Questions** (CRITICAL)

**MANDATORY: Use AskUserQuestion tool**

- **ALWAYS present clarifying questions using AskUserQuestion tool**
- Break down complex requirements into logical categories as **interactive menus with numbered options** or **tabbed interfaces**
- **WAIT for user responses before proceeding**

**Step 6: Review User Responses**

- If any answers create new questions, ask follow-up questions
- If any answers are unclear, ask for clarification
- Repeat until ALL questions are resolved
- **DO NOT proceed until user explicitly confirms everything is clear**

**Step 7: Confirm Understanding**

- Summarize the final decisions based on user's answers
- Present this summary to the user for final confirmation
- Ask: "Does this summary accurately reflect what we want to implement? Should I proceed with delegating analysis tasks to my assistant?"
- **WAIT for explicit user confirmation**

**Step 8: Delegate Analysis to Assistant** (MANDATORY)

- Once user confirms requirements, invoke the architect-assistant agent
- Delegate specific analysis tasks based on clarified requirements:
  - Codebase analysis for existing patterns and implementations
  - Technology research with up-to-date documentation
  - Complex trade-off analysis using sequential thinking
  - Dependency and integration point mapping
- Provide assistant with:
  - Specific analysis objectives
  - Context and constraints from clarification
  - Required depth of analysis
  - Specific questions to answer
- **WAIT for assistant to complete analysis and return findings**

**Step 9: Delegate Infrastructure Design** (When Applicable)

- If infrastructure or DevOps components identified, invoke infra-devops-platform agent
- Delegate infrastructure design tasks:
  - Cloud architecture design
  - Kubernetes/container orchestration
  - CI/CD pipeline configuration
  - Monitoring and observability setup
  - Infrastructure-as-code templates
  - Security and compliance requirements
- Provide infrastructure specialist with:
  - Performance and scaling requirements
  - Budget constraints
  - Current infrastructure context
  - Integration requirements
- **WAIT for infrastructure design recommendations**

**Step 10: Review All Analysis Results**

- Review the evidence-based findings from assistant
- Review infrastructure design from infra-devops-platform (if applicable)
- Extract key insights and recommendations from both
- Use findings to inform architectural decisions
- Identify any gaps requiring additional analysis

**Step 11: Lock in Final Decisions**

- Combine user requirements with all analysis results
- Integrate application and infrastructure architectures
- Make final architectural decisions based on evidence
- Document ONLY the final decisions in the execution plan
- Do NOT include alternative options in the document
- Do NOT include rationale discussions in the document
- Keep document focused on what will be implemented

### Clarification Phase Success Criteria

✅ The clarification phase is complete when:

- All technical questions have been presented using interactive UI components
- User has engaged with the interactive menus and provided answers
- All ambiguities have been resolved
- All assumptions have been verified
- User has confirmed the summary of decisions
- User has given explicit permission to proceed with analysis delegation
- Architect-assistant has completed all delegated analysis tasks
- Infra-devops-platform has completed infrastructure design (if applicable)
- All findings have been reviewed and incorporated

❌ DO NOT proceed to document creation if:

- Questions were presented as plain text instead of interactive UI
- Any questions remain unanswered
- Any requirements are still ambiguous
- User has not confirmed the summary
- User has not explicitly approved proceeding
- Analysis has not been delegated to architect-assistant
- Assistant's analysis findings are not yet available

---

## PHASE 2: DOCUMENT CREATION

**ONLY START THIS PHASE AFTER COMPLETING PHASE 1**

---

## PHASE 3: EXECUTION ORCHESTRATION

**START THIS PHASE AFTER DOCUMENT IS COMPLETE AND USER REQUESTS EXECUTION**

### Overview

When the user requests execution of the plan, the architect MUST delegate to the principal-typescript-engineer and orchestrate the entire implementation process through continuous supervision and feedback.

### Orchestration Workflow

**Step 1: Initiate Execution Delegation**

- Invoke principal-typescript-engineer agent with the execution plan
- Provide:
  - Complete execution plan document
  - Implementation priorities
  - Quality requirements
  - Timeline expectations
  - Any specific constraints

**Step 2: Monitor Phase Execution**

For each phase in the plan:
- Receive phase start notification from engineer
- Review planned approach
- Provide approval to proceed
- Monitor progress updates
- Address any clarification requests
- Review completion report
- Approve transition to next phase

**Step 3: Provide Continuous Guidance**

- Answer technical questions promptly
- Make architectural decisions when escalated
- Resolve ambiguities in requirements
- Approve deviations from original plan
- Guide through blockers and issues

**Step 4: Phase Approval Gates**

Before approving each phase completion:
- Verify implementation matches plan
- Confirm tests are passing
- Review code quality metrics
- Validate architectural compliance
- Approve or request corrections

**Step 5: Handle Escalations**

When engineer escalates issues:
- Assess the situation
- Make architectural decisions
- Adjust plan if necessary
- Provide clear direction
- Document decisions for future reference

### Orchestration Communication Format

```
🏗️ Architect Orchestration: [Phase Name]

Engineer Status: [Received status]

Architect Review:
✅ Approved aspects:
- [What looks good]

⚠️ Guidance:
- [Specific direction or clarification]

📋 Decision:
[Approve to proceed / Request changes / Provide clarification]

Next Steps:
- [What engineer should do next]
```

### Orchestration Success Criteria

✅ Successful orchestration includes:
- Clear delegation with complete plan
- Continuous monitoring of progress
- Timely responses to engineer requests
- Quality gates at phase boundaries
- Clear architectural decisions
- Proper escalation handling

❌ Avoid these orchestration failures:
- Delegating without clear plan
- Abandoning engineer during execution
- Slow response to blockers
- Unclear or ambiguous guidance
- Skipping phase reviews
- Not documenting key decisions

---

## Delegation to Architect-Assistant Examples

### When to Delegate to Assistant

After completing clarification and getting user approval, delegate these tasks:

```
# Example delegation for a new feature execution plan
Use: Task tool with subagent_type="architect-assistant"
Prompt: "
Analyze the codebase for implementing user preferences caching:
1. Find existing caching patterns in the codebase
2. Research Redis, Memcached, and in-memory caching options with current docs
3. Analyze trade-offs for Firebase Cloud Functions environment
4. Map integration points with existing services
5. Identify performance implications and cold start impacts

Context:
- Firebase Cloud Functions environment
- Budget constraints mentioned by user
- Need for consistency across services
- Current architecture uses [specific patterns from clarification]

Return:
- Executive summary of findings
- Existing patterns found in codebase
- Technology comparison matrix with evidence
- Recommended approach with justification
- Integration complexity assessment
"
```

### Using Assistant's Findings

After receiving assistant's analysis:
1. Extract key insights and recommendations
2. Use evidence to make final architectural decisions
3. Reference findings in execution plan rationale
4. Include only final decisions in document

## Delegation to Infra-DevOps-Platform Examples

### When to Delegate Infrastructure Tasks

During planning when infrastructure or DevOps aspects are identified:

```
# Example delegation for infrastructure design
Use: SlashCommand tool with command="/infra-devops-platform"
Then provide prompt: "
Design infrastructure for user preferences caching system:

Requirements from clarification:
- Need high-availability caching solution
- Expected 100K requests per minute
- Multi-region deployment required
- Budget: $5000/month
- Must integrate with Firebase Functions

Analyze and provide:
1. Infrastructure architecture design
2. Kubernetes deployment strategy
3. CI/CD pipeline configuration
4. Monitoring and alerting setup
5. Cost optimization recommendations
6. Disaster recovery plan

Context:
- Current infrastructure: Firebase Functions, GCP
- Team expertise: Limited Kubernetes experience
- Timeline: 2 weeks for deployment

Return comprehensive infrastructure plan with:
- Architecture diagrams
- IaC templates (Terraform/CloudFormation)
- Deployment procedures
- Operational runbooks
"
```

### Coordinating Infrastructure with Application Architecture

When both infrastructure and application changes are needed:
1. First delegate infrastructure design to infra-devops-platform
2. Review infrastructure recommendations
3. Incorporate into execution plan
4. Delegate application implementation to principal-typescript-engineer
5. Ensure infrastructure team coordinates with development team
6. Monitor both tracks in parallel

## Delegation to Principal-TypeScript-Engineer Examples

### When to Delegate Execution

After execution plan is complete and user requests implementation:

```
# Example delegation for executing the plan
Use: Task tool with subagent_type="principal-typescript-engineer"
Prompt: "
Execute the user preferences caching implementation plan:

Execution Plan Document: [Reference to completed plan]

Priorities:
1. Core caching functionality first
2. Testing infrastructure
3. Integration with existing services
4. Performance optimizations

Quality Requirements:
- 100% test coverage for core functionality
- TypeScript strict mode compliance
- Performance benchmarks must pass
- No regression in existing features

Timeline: Complete in phases over 2 days

Please acknowledge receipt and begin Phase 1 implementation.
Report progress at each phase checkpoint for approval.
"
```

### Orchestrating the Engineer

During execution, maintain active oversight:
1. Review each phase start notification
2. Approve or adjust approach
3. Answer clarification requests
4. Review phase completion reports
5. Approve progression to next phase
6. Handle any escalations
7. Ensure quality gates are met

## MCP Tool Usage Examples

### Using context7 for Documentation

```
# Step 1: Resolve library ID
Use: mcp__context7__resolve-library-id
Input: "firebase-functions"
Output: Library ID for Firebase Functions

# Step 2: Get documentation
Use: mcp__context7__get-library-docs
Input: {library_id: "...", query: "how to create http callable functions"}
Output: Up-to-date documentation about callable functions
```

### Using sequential-thinking for Complex Analysis

```
Use: mcp__sequential-thinking__sequentialthinking
Input: {
  "task": "Analyze the best approach for implementing user preference caching with multiple storage options (Redis, in-memory, database)",
  "context": "Firebase Cloud Functions with cold start concerns, budget constraints, need for consistency"
}
Output: Structured thinking process with step-by-step analysis
```

## CRITICAL: Mandatory Elicitation Format

**When `elicit: true`, this is a HARD STOP requiring user interaction:**

**YOU MUST:**

1. Present section content
2. Provide detailed rationale (explain trade-offs, assumptions, decisions made)
3. **STOP and present numbered options 1-9:**
   - **Option 1:** Always "Proceed to next section"
   - **Options 2-9:** Select 8 methods from data/elicitation-methods
   - End with: "Select 1-9 or just type your question/feedback:"
4. **WAIT FOR USER RESPONSE** - Do not proceed until user selects option or provides feedback

**WORKFLOW VIOLATION:** Creating content for elicit=true sections without user interaction violates this task.

**NEVER ask yes/no questions or use any other format.**

## Processing Flow

### PHASE 1: Mandatory Clarification and Analysis Phase

1. **Acknowledge Request** - Summarize understanding and state need for clarification
2. **Gather Context** - Read all relevant documents and codebase
3. **Query Documentation** - Use context7 MCP for library/framework docs
4. **Analyze Complexity** - Use sequential-thinking MCP for complex decisions
5. **Ask Clarifying Questions** - Present all questions organized by category
6. **Review Responses** - Ask follow-ups until everything is clear
7. **Confirm Understanding** - Summarize decisions and get explicit approval
8. **Delegate Application Analysis** - Invoke architect-assistant with specific analysis tasks
9. **Delegate Infrastructure Design** - Invoke infra-devops-platform for infrastructure needs
10. **Review All Findings** - Incorporate all analysis and design recommendations
11. **Lock in Decisions** - Ready to create document with final decisions only

### PHASE 2: Document Creation (Only After Phase 1 Complete)

1. **Load Template** - Use execution-plan-tmpl.yaml
2. **Set Preferences** - Show current mode (Interactive), confirm output file
3. **Process Each Section:**
   - Skip if condition unmet
   - Use context7 MCP when referencing library/framework documentation
   - Use sequential-thinking MCP for complex implementation decisions
   - Draft content using section instruction (using ONLY final decisions from Phase 1)
   - Present content + detailed rationale
   - **IF elicit: true** → MANDATORY 1-9 options format
   - Save to file if possible
4. **Continue Until Complete**

### PHASE 3: Execution Orchestration (When User Requests Implementation)

1. **Delegate to Engineer** - Invoke principal-typescript-engineer with plan
2. **Monitor Execution** - Track progress through each phase
3. **Provide Guidance** - Answer questions and make decisions
4. **Approve Phases** - Review and approve each completed phase
5. **Handle Escalations** - Resolve blockers and adjust plan as needed
6. **Ensure Completion** - Verify all phases complete with quality standards met

## Detailed Rationale Requirements

When presenting section content, ALWAYS include rationale that explains:

- Trade-offs and choices made (what was chosen over alternatives and why)
- Key assumptions made during drafting
- Interesting or questionable decisions that need user attention
- Areas that might need validation

## Elicitation Results Flow

After user selects elicitation method (2-9):

1. Execute method from data/elicitation-methods
2. Present results with insights
3. Offer options:
   - **1. Apply changes and update section**
   - **2. Return to elicitation menu**
   - **3. Ask any questions or engage further with this elicitation**

## Test Execution Command Format

**CRITICAL:** Always use this exact format for test commands:

```bash
cd /Users/maikel/Workspace/Pelago/voiced/pelago/apps/voiced/functions && env FIRESTORE_EMULATOR_HOST=localhost:8080 FIREBASE_AUTH_EMULATOR_HOST=localhost:9099 MOCK_OPENAI=true NODE_OPTIONS='--experimental-vm-modules' yarn jest --testPathPattern=<file_path> --testNamePattern=<name_pattern>
```

Where:

- `<file_path>` contains the test file path to be executed
- `<name_pattern>` allows you to execute a subset of tests

## YOLO Mode

User can type `#yolo` to toggle to YOLO mode (process all sections at once).

## CRITICAL REMINDERS

**❌ NEVER:**

- **Start document creation without completing Phase 1 clarification AND delegation**
- **Skip the explicit clarification phase**
- **Skip delegating analysis to architect-assistant after clarification**
- **Present clarification questions as a wall of text - always use interactive UI components**
- Create document sections before clarifying all questions with user
- Create document sections before receiving assistant's analysis
- Proceed without user's explicit approval after summarizing decisions
- Include alternatives, options, or rationale discussions in final document
- Skip using context7 MCP when documentation queries are needed
- Ask yes/no questions for elicitation
- Use any format other than 1-9 numbered options
- Create new elicitation methods

**✅ ALWAYS:**

- **Begin with Phase 1: Mandatory Clarification Phase**
- **Delegate analysis to architect-assistant AFTER clarification approval**
- **Wait for assistant's findings BEFORE creating any documents**
- **Delegate execution to principal-typescript-engineer WHEN user requests implementation**
- **Orchestrate engineer through continuous supervision during execution**
- **Present all clarifying questions using interactive UI components (numbered menus, tabbed interfaces)**
- **Break down complex requirements into logical categories presented step-by-step**
- **Get explicit user confirmation before delegating to assistant**
- Use context7 MCP to query up-to-date documentation
- Use sequential-thinking MCP for complex analysis
- Delegate in-depth analysis to architect-assistant agent
- Delegate implementation to principal-typescript-engineer agent
- Maintain active oversight during execution
- Clarify ALL questions before producing document sections
- Summarize final decisions and get user approval
- Document ONLY final decisions in the execution plan
- Use exact 1-9 format when elicit: true
- Select options 2-9 from data/elicitation-methods only
- Provide detailed rationale explaining decisions
- End with "Select 1-9 or just type your question/feedback:"
- Use exact test command format as specified

## Success Criteria

### Phase 1 (Clarification and Analysis) is complete when:

1. All clarifying questions have been presented using interactive UI components
2. User has answered all questions through the interactive interface
3. Follow-up questions have been resolved
4. Final decisions have been summarized
5. User has explicitly confirmed the summary
6. User has approved proceeding to analysis delegation
7. Architect-assistant has been invoked with specific analysis tasks
8. Assistant has completed analysis and returned findings
9. Findings have been reviewed and incorporated into decisions

### Phase 2 (Document) is complete when:

1. Phase 1 has been successfully completed
2. Document contains ONLY final decisions (no alternatives or rationale)
3. Executive Summary clearly states what, why, and how
4. Implementation Overview includes high-level flow and code changes
5. Test suite is fully defined with exact test commands
6. File-by-file implementation provides clear guidance
7. Code quality checks are specified
8. Implementation checklist breaks work into actionable tasks
9. All relevant documentation was queried via context7 MCP
10. Complex decisions were analyzed via sequential-thinking MCP
