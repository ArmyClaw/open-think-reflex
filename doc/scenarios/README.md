# Open-Think-Reflex Scenarios

> **Version**: v2.1  
> **Status**: Draft  
> **Scope**: CLI Use Case Design for v1.0 (AI Input Accelerator)

---

## Table of Contents

1. [Overview](#1-overview)
2. [User Personas](#2-user-personas)
3. [Core Workflow Scenarios](#3-core-workflow-scenarios)
4. [Thinking & Writing Scenarios](#4-thinking--writing-scenarios)
5. [Analysis & Decision Scenarios](#5-analysis--decision-scenarios)
6. [Problem Solving Scenarios](#6-problem-solving-scenarios)
7. [Learning & Research Scenarios](#7-learning--research-scenarios)
8. [Documentation Scenarios](#8-documentation-scenarios)
9. [Edge Cases](#9-edge-cases)
10. [Scenario Comparison Matrix](#10-scenario-comparison-matrix)

---

## 1. Overview

### 1.1 Purpose

This document provides concrete CLI use case scenarios for Open-Think-Reflex, demonstrating how the reflex system accelerates human thinking and AI collaboration.

### 1.2 Core Concept

```
Human Thinking                    AI Collaboration
─────────────────────────────────────────────────────────────────
思考过程                    AI辅助生成
   ↓                              ↓
反复使用 → 形成反射           快速触发生成
   ↓                              ↓
自动化响应                    加速思考过程
```

**核心观点**：
- 反射 = 思维模式的快捷方式
- 重复强化 → 形成自动化响应
- 不用衰减 → 保持思维活跃

### 1.3 Three-Layer Interaction Design

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        三层交互界面设计                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  上层：思维链树展示层（从左到右展开）                                │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐            │
│  │ 用户API  │ ───► │ 分页    │ ───► │ pageSize│            │
│  │  85%    │      │  85%    │      │ pageNum │            │
│  └─────────┘      └─────────┘      └─────────┘      │
│                                                                      │
│  视觉反馈：                                                          │
│  ├── 已选分支：高亮显示                                              │
│  ├── 未选分支：灰色/淡化                                             │
│  └── 选中态：边框闪烁                                                 │
│                                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  中层：当前输出结果层                                                │
│  ┌─────────────────────────────────────────────────────────────┐     │
│  │                                                              │     │
│  │              AI输出的结果内容...                                │     │
│  │                                                              │     │
│  └─────────────────────────────────────────────────────────────┘     │
│                                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  下层：用户输入层                                                    │
│  ┌─────────────────────────────────────────────────────────────┐     │
│  │ > 用户输入命令或问题...                            [Enter] │     │
│  └─────────────────────────────────────────────────────────────┘     │
│                                                                      │
│  操作提示：[↑↓]选择分支 [→/Enter]确认展开 [←/Esc]返回 [空格]触发生成│
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

> **注意**：本文档中的ASCII表格仅为展示分层结构，实际CLI交互中采用上中下三层布局。

### 1.3 Design Principles

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Reflex System Design Principles                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. Terminal-First                                                      │
│     - All interactions via CLI                                          │
│     - Keyboard shortcuts for efficiency                                 │
│     - Clear visual feedback in terminal                                 │
│                                                                      │
│  2. Thought Reflection                                                 │
│     - Capture recurring thinking patterns                               │
│     - Accelerate repetitive thought processes                          │
│     - Enable progressive deepening of insights                          │
│                                                                      │
│  3. Adaptive Learning                                                  │
│     - System learns from usage patterns                                │
│     - Strength increases with repetition                              │
│     - Decay prevents stale thinking                                    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.4 Application Domains

| Domain | Examples | User Types |
|--------|----------|-------------|
| **Thinking & Writing** | Articles, outlines, content structures | Writers, researchers |
| **Analysis & Decision** | SWOT, root cause, comparative analysis | Analysts, managers |
| **Problem Solving** | Debugging frameworks, troubleshooting | All users |
| **Learning & Research** | Concept explanation, comparative study | Students, learners |
| **Documentation** | Process docs, project overviews | All users |

---

## 2. User Personas

### 2.1 Persona Profiles

#### Persona 1: Content Strategist (Alex)

```
================================================================================
Name: Alex
Role: Content Strategist
Experience: 6 years

Work Pattern:
  - Creates 25+ content pieces/month
  - Uses consistent frameworks
  - Values clarity and efficiency

Goals:
  1. Faster content creation workflow
  2. Consistent structure across content
  3. Quick reference to frameworks

Pain Points:
  - Repeating content structures
  - Hard to recall analytical frameworks
  - Context switching between topics

Reflex Candidate Examples:
  - "article-outline" → generates article structure
  - "swot-framework" → generates SWOT analysis template
  - "competitive-analysis" → generates comparison structure
================================================================================
```

#### Persona 2: Business Analyst (Sam)

```
================================================================================
Name: Sam
Role: Business Analyst
Experience: 8 years

Work Pattern:
  - Creates 20+ analysis reports/month
  - Uses consistent analytical frameworks
  - Values accuracy over speed

Goals:
  1. Faster analysis workflow
  2. Consistent report structures
  3. Easy reference to methodologies

Pain Points:
  - Repeating analytical steps
  - Hard to maintain consistent frameworks
  - Knowledge gaps in new domains

Reflex Candidate Examples:
  - "root-cause" → generates problem analysis framework
  - "decision-matrix" → generates comparison template
  - "risk-assessment" → generates risk matrix
================================================================================
```

#### Persona 3: Continuous Learner (Jordan)

```
================================================================================
Name: Jordan
Role: Continuous Learner
Experience: Self-taught professional

Work Pattern:
  - Studies 5+ topics/month
  - Creates study notes and summaries
  - Values understanding over speed

Goals:
  1. Faster knowledge acquisition
  2. Easy concept review
  3. Progressive learning paths

Pain Points:
  - Hard to retain information
  - Repeating learning paths
  - Difficulty connecting concepts

Reflex Candidate Examples:
  - "concept-summary" → generates concept summary
  - "compare-technologies" → generates comparison table
  - "study-guide" → generates learning plan
================================================================================
```

### 2.2 Persona Usage Patterns

| Persona | Daily Reflections | Weekly New | Typical Session Length |
|---------|-------------------|------------|----------------------|
| Alex | 10-20 | 2-3 | 5-10 min bursts |
| Sam | 5-15 | 1-2 | 15-30 min sessions |
| Jordan | 5-10 | 3-5 | 20-40 min sessions |

---

## 3. Core Workflow Scenarios

### 3.1 Scenario: First-Time User Onboarding

#### Description

New user installs Open-Think-Reflex, configures AI provider, and creates first reflex through guided interaction.

#### User Story

```
As a new user
I want to configure the system and create my first reflex
So that I can experience the speed of reflex-based thinking
```

#### Workflow

```bash
# Step 1: Installation
$ npm install -g @openclaw/reflex

# Step 2: Initial Setup
$ otr init

# Interactive Configuration:
# ? Select AI provider: [Use arrow keys]
#    > claude (Recommended)
#      openai
#      local
# ? Enter your API key: **********
# ? Select data directory: ~/.otr
# ? Configure matching strategies: y
#     - exact: enabled (default)
#     - keyword: enabled (default)
#     - semantic: disabled (default)
# ? Configure decay settings: y
#     - short: 7 days (default)
#     - medium: 70 days (default)
#     - long: 700 days (default)
# ? Enable automatic backup: y

# Step 3: First Interaction
$ otr "How to analyze market trends"

# Output:
# No matching reflex found.
# [Space] to generate with AI, or [Esc] to cancel.
#
# [Space] Generating with Claude...
#
# ┌─────────────────────────────────────────┐
# │ Thinking Chain                          │
# │ [root] ──► [Branch A] ──► [Branch A1]   │
# │                                         │
# └─────────────────────────────────────────┘
# ┌─────────────────────────────────────────┐
# │ Output                                   │
# │ # Market Trend Analysis                 │
# │ ## 1. Data Collection                  │
# │ ## 2. Pattern Identification           │
# │ ## 3. Trend Extrapolation              │
# │ ## 4. Risk Assessment                  │
# │                                         │
# └─────────────────────────────────────────┘
# ┌─────────────────────────────────────────┐
# │ Input:                                  │
# │ > _                                     │
# └─────────────────────────────────────────┘

# Step 4: Save as Reflex (Optional)
# ? Save this pattern as a reflex? (y/n)
# y
# ? Enter trigger keyword: market-analysis
# ? Add tags: [analysis, business, strategy]
# Reflex "market-analysis" created with strength 10%
```

#### Expected Outcomes

```
✓ User completes configuration
✓ User experiences first AI generation
✓ Optional: User creates first reflex
✓ System records usage metrics
```

### 3.2 Scenario: Rapid Reflex Execution

#### Description

Experienced user quickly triggers an existing reflex for a frequently used thinking pattern.

#### User Story

```
As an experienced user
I want to trigger a high-strength reflex with minimal keystrokes
So that I can complete repetitive thinking tasks in seconds
```

#### Workflow

```bash
# User has a "market-analysis" reflex with strength 95%

$ otr market-analysis

# Output:
# ┌─────────────────────────────────────────┐
# │ Thinking Chain                          │
# │ [market-analysis] ──► [Data Collection] ✓    │
# │ [market-analysis] ──► [Pattern ID] ◉          │
# │ [market-analysis] ──► [Risk Assessment]       │
# │ [market-analysis] ──► [Recommendations]      │
# │                                         │
# │ Selected: [Pattern ID] (95% confidence)        │
# └─────────────────────────────────────────┘
# ┌─────────────────────────────────────────┐
# │ Output: Generated Framework               │
# │ ## 1. Data Sources                     │
# │    - Sales data                        │
# │    - Customer feedback                  │
# │    - Competitor analysis                │
# │    - Market research                   │
# │                                         │
# │ ## 2. Pattern Identification           │
# │    - Trend direction (up/down/stable)  │
# │    - Growth rate calculation           │
# │    - Seasonality analysis              │
# │                                         │
# │ ## 3. Key Insights                    │
# │    - Primary drivers                  │
# │    - Emerging opportunities            │
# │    - Potential threats                 │
# └─────────────────────────────────────────┘
# ┌─────────────────────────────────────────┐
# │ Input:                                  │
# │ > _ [↓/↑ select, → expand, Space exec] │
# └─────────────────────────────────────────┘

# User wants full framework, presses [→] to expand
# [→] Expand to show all sections
# [Space] Execute and copy to clipboard

# System output:
# ✓ Content copied to clipboard
# ✓ Reflex "market-analysis" reinforced (strength: 95% → 98%)
```

#### Expected Outcomes

```
✓ Fast reflex triggering (< 1 second)
✓ Visual confirmation of reflex match
✓ One-key execution with [Space]
✓ Automatic reinforcement of reflex strength
```

### 3.3 Scenario: Reflex Strength Progression

#### Description

User observes how reflex strength increases with repeated use and eventually becomes "deep reflex."

#### User Story

```
As a user
I want to see my frequently used patterns become stronger reflexes
So that I can trust the system to remember my preferred thinking patterns
```

#### Workflow

```bash
# Day 1: First use
$ otr "market-analysis"
# Strength: 10% (New impression)

# Day 2: Second use
$ otr "market-analysis"
# Strength: 30% (Reinforced)

# Day 3: Third use
$ otr "market-analysis"
# Strength: 50% (Threshold reached - now activated automatically)

# Day 7: Seventh use
$ otr "market-analysis"
# Strength: 85% (Strong reflex)

# Day 14: Fourteenth use
$ otr "market-analysis"
# Strength: 100% (Deep reflex - permanent)

# User checks reflex status
$ otr pattern show market-analysis

# Output:
# ┌─────────────────────────────────────────┐
# │ Pattern: market-analysis                   │
# │ Status: DEEP REFLEX (Permanent)          │
# │ Strength: 100%                          │
# │ Threshold: 50%                           │
# │ Reinforcement Count: 14                  │
# │ Decay Count: 0                           │
# │ Created: 2024-01-01                     │
# │ Last Used: 2024-01-14                    │
# │ Tags: [analysis, business, strategy]    │
# │                                          │
# │ Response:                               │
# │ # Market Trend Analysis Framework         │
# │ ...                                      │
# └─────────────────────────────────────────┘
```

#### Expected Outcomes

```
✓ Visual strength indicator in output
✓ Threshold notification at 50%
✓ Permanent status at 100%
✓ Reflex history tracking
```

---

## 4. Thinking & Writing Scenarios

### 4.1 Scenario: Article Outline Generation

#### Description

Content creator generates a structured outline for a new article.

#### User Story

```
As a content creator
I want to generate a consistent article structure
So that I can maintain writing consistency
```

#### Workflow

```bash
# User triggers outline generation
$ otr "Write article about AI trends"

# System shows matching reflexes:
# ┌─────────────────────────────────────────┐
# │ Thinking Chain                          │
# │ [article] ──► [tech] ──► [tutorial] ◉ │
# │ [article] ──► [tech] ──► [news-analysis]│
# │ [article] ──► [business] ──► [case-study]│
# │ [article] ──► [general] ──► [how-to] │
# │                                         │
# │ Confidence: 85%                          │
# └─────────────────────────────────────────┘

# User selects "news-analysis" branch
# [↓] Select "news-analysis"
# [Space] Generate

# Output:
# ┌─────────────────────────────────────────┐
# │ Generated Outline                         │
# │ =====================================  │
# │                                          │
# │ # AI Trends 2024: What's Next?        │
# │                                          │
# │ ## 1. Executive Summary               │
# │    - Key findings (3-4 sentences)      │
# │    - Main thesis statement             │
# │    - Reader take-away                  │
# │                                          │
# │ ## 2. Current Landscape               │
# │    - Market size and growth           │
# │    - Key players overview              │
# │    - Technology adoption trends        │
# │                                          │
# │ ## 3. Emerging Trends                 │
# │    - Trend 1: [Specific trend]        │
# │      - Evidence and data               │
# │      - Industry impact                 │
# │      - Future prediction              │
# │    - Trend 2: [Specific trend]        │
# │    - Trend 3: [Specific trend]        │
# │                                          │
# │ ## 4. Challenges                      │
# │    - Technical barriers               │
# │    - Ethical concerns                  │
# │    - Regulatory landscape              │
# │                                          │
# │ ## 5. Expert Opinions                │
# │    - Quote 1: [Industry expert]       │
# │    - Quote 2: [Researcher]          │
# │                                          │
# │ ## 6. Future Outlook                  │
# │    - 1-year prediction                │
# │    - 5-year projection                 │
# │    - Call to action                   │
# │                                          │
# │ ## 7. References                      │
# │    - Key sources                      │
# │    - Related articles                  │
# │    - Further reading                  │
# └─────────────────────────────────────────┘

# User actions:
# [Enter] Copy to clipboard
# [n] Don't save as new reflex (already exists)
```

#### Expected Outcomes

```
✓ Complete, well-structured outline
✓ Consistent formatting
✓ Ready for expansion
✓ One-click copy to clipboard
```

### 4.2 Scenario: Content Framework

#### Description

Writer generates consistent content frameworks for different content types.

#### User Story

```
As a content strategist
I want to generate content frameworks
So that I can maintain consistency across all content
```

#### Workflow

```bash
# User triggers framework generation
$ otr "Create product launch announcement"

# System shows:
# ┌─────────────────────────────────────────┐
# │ Thinking Chain                          │
# │ [content] ──► [product] ──► [launch] ◉ │
# │ [content] ──► [product] ──► [update] │
# │ [content] ──► [product] ──► [legacy] │
# │ [content] ──► [email] ──► [newsletter] │
# │                                         │
# │ Confidence: 78%                          │
# └─────────────────────────────────────────┘

# User selects "launch" branch
# [Space] Generate

# Output:
# ┌─────────────────────────────────────────┐
# │ Generated Framework                      │
# │ =====================================  │
# │                                          │
# │ # Product Launch Announcement          │
# │                                          │
# │ ## 1. Hook (Attention Grabber)        │
# │    - Surprising statistic or fact       │
# │    - Problem statement                 │
# │    - Bold promise                      │
# │                                          │
# │ ## 2. Problem Agitation                │
# │    - Current pain points               │
# │    - Cost of inaction                  │
# │    - Emotional connection              │
# │                                          │
# │ ## 3. Solution Introduction           │
# │    - Product name reveal               │
# │    - One-line value proposition       │
# │    - Core benefit highlighted          │
# │                                          │
# │ ## 4. Features & Benefits              │
# │    - Feature 1 → Benefit 1            │
# │    - Feature 2 → Benefit 2            │
# │    - Feature 3 → Benefit 3            │
# │                                          │
# │ ## 5. Social Proof                    │
# │    - Beta tester quotes               │
# │    - Early adoption metrics            │
# │    - Expert endorsements               │
# │                                          │
# │ ## 6. Offer Details                    │
# │    - Launch date                      │
# │    - Pricing structure                 │
# │    - Early bird discount              │
# │                                          │
# │ ## 7. Call to Action                  │
# │    - Primary CTA                      │
# │    - Secondary CTA                    │
# │    - Urgency element                  │
# │                                          │
# │ ## 8. Closing                         │
# │    - Vision statement                 │
# │    - Invitation to join               │
# │    - Contact information              │
# └─────────────────────────────────────────┘

# [Enter] Copy to clipboard
# [y] Save reflex "product-launch-framework"
```

#### Expected Outcomes

```
✓ Comprehensive framework
✓ Conversion-focused structure
✓ Easy to customize
✓ Ready for content creation
```

---

## 5. Analysis & Decision Scenarios

### 5.1 Scenario: SWOT Analysis Framework

#### Description

Analyst generates a SWOT analysis framework for strategic planning.

#### User Story

```
As a business analyst
I want to generate a SWOT analysis structure
So that I can consistently evaluate business situations
```

#### Workflow

```bash
# User triggers SWOT analysis
$ otr "SWOT analysis for our new product"

# System shows analysis options:
# ┌─────────────────────────────────────────┐
# │ Thinking Chain                          │
# │ [swot] ──► [product] ──► [tech] ◉   │
# │ [swot] ──► [product] ──► [service] │
# │ [swot] ──► [company] ──► [startup] │
# │ [swot] ──► [market] ──► [competitor]│
# │                                         │
# │ Context: New SaaS product launch       │
# └─────────────────────────────────────────┘

# User selects "tech product" branch
# [Space] Generate

# Output:
# ┌─────────────────────────────────────────┐
# │ Generated SWOT Analysis                  │
# │ =====================================  │
# │                                          │
# │ # SWOT Analysis: [Product Name]        │
# │                                          │
# │ ## 1. STRENGTHS                       │
# │                                          │
# │ ### Internal Factors                    │
# │ - Technology advantage: [Your USP]    │
# │ - Team expertise: [Core competencies] │
# │ - Resource availability: [Funding/IP] │
# │ - Process efficiency: [Methodology]    │
# │ - Brand recognition: [Market presence]│
# │                                          │
# │ ### Questions to Consider              │
# │ - What do we do better than anyone?   │
# │ - What unique resources can we use?    │
# │ - What do others see as our strengths?│
# │                                          │
# │ ## 2. WEAKNESSES                      │
# │                                          │
# │ ### Internal Factors                    │
# │ - Limited resources: [Budget/People]  │
# │ - Technology gaps: [Missing features] │
# │ - Brand weakness: [Market awareness]  │
# │ - Knowledge gaps: [Experience level]  │
# │ - Process issues: [Inefficiencies]     │
# │                                          │
# │ ### Questions to Consider              │
# │ - What could we improve?               │
# │ - What do competitors do better?       │
# │ - What resources are missing?          │
# │                                          │
# │ ## 3. OPPORTUNITIES                   │
# │                                          │
# │ ### External Factors                    │
# │ - Market trends: [Growing demand]     │
# │ - Technology changes: [New platforms] │
# │ - Regulatory changes: [New policies]  │
# │ - Partnership possibilities: [Alliances]│
# │ - Competitor weaknesses: [Market gaps] │
# │                                          │
# │ ### Questions to Consider              │
# │ - What market trends favor us?        │
# │ - What technologies can we leverage?   │
# │ - What gaps exist in the market?      │
# │                                          │
# │ ## 4. THREATS                         │
# │                                          │
# │ ### External Factors                    │
# │ - Competition: [Competitor actions]   │
# │ - Market changes: [Shifting demands]  │
# │ - Technology changes: [Disruption]    │
# │ - Regulatory risks: [Compliance]       │
# │ - Economic factors: [Budget concerns]  │
# │                                          │
# │ ### Questions to Consider              │
# │ - What threats does our weakness create?│
# │ - What are competitors doing?         │
# │ - What market changes could hurt us?   │
# │                                          │
# │ ## 5. ACTION ITEMS                   │
# │                                          │
# │ ### Leverage Strengths                 │
# │ - [Action 1]                          │
# │ - [Action 2]                          │
# │                                          │
# │ ### Address Weaknesses                │
# │ - [Action 1]                          │
# │ - [Action 2]                          │
# │                                          │
# │ ### Capture Opportunities             │
# │ - [Action 1]                          │
# │ - [Action 2]                          │
# │                                          │
# │ ### Mitigate Threats                  │
# │ - [Action 1]                          │
# │ - [Action 2]                          │
# └─────────────────────────────────────────┘

# User actions:
# [Enter] Copy to clipboard
# [y] Save reflex "swot-product-analysis"
```

#### Expected Outcomes

```
✓ Complete SWOT framework
✓ Questions to guide analysis
✓ Action items section
✓ Easy to customize
```

### 5.2 Scenario: Decision Matrix

#### Description

Analyst generates a decision matrix framework for comparing options.

#### User Story

```
As a decision maker
I want to compare multiple options systematically
So that I can make data-driven decisions
```

#### Workflow

```bash
# User needs to compare options
$ otr "Compare vendors for our new tool"

# System shows:
# ┌─────────────────────────────────────────┐
# │ Thinking Chain                          │
# │ [decision] ──► [vendor] ──► [tech] ◉ │
# │ [decision] ──► [vendor] ──► [service]│
# │ [decision] ──► [solution] ──► [build-vs-buy]│
# │ [decision] ──► [solution] ──► [platform] │
# │                                         │
# │ Context: Comparing 4 vendors           │
# └─────────────────────────────────────────┘

# User selects "tech vendors" branch
# [Space] Generate

# Output:
# ┌─────────────────────────────────────────┐
# │ Generated Decision Matrix                │
# │ =====================================  │
# │                                          │
# │ # Vendor Comparison Matrix              │
# │                                          │
# │ ## Evaluation Criteria (Weight: 1-5)   │
# │                                          │
# │ | Criteria | Weight | Description |       │
# │ |----------|--------|-------------|       │
# │ | Cost | 5 | Total cost of ownership |  │
# │ | Features | 4 | Feature completeness | │
# │ | Usability | 4 | Ease of use |       │
# │ | Support | 3 | Quality of support |    │
# │ | Scalability | 4 | Growth potential |  │
# │ | Security | 5 | Security posture |    │
# │ | Integration | 3 | API availability |  │
# │ | Reputation | 3 | Market presence |    │
# │                                          │
# │ ## Vendor Ratings (1-10)               │
# │                                          │
# │ | Criteria | Vendor A | Vendor B | Vendor C | Vendor D |      │
# │ |----------|----------|----------|----------|----------|      │
# │ | Cost | 7 | 8 | 6 | 5 |               │
# │ | Features | 8 | 7 | 9 | 6 |          │
# │ | Usability | 6 | 9 | 7 | 8 |          │
# │ | Support | 7 | 6 | 8 | 5 |           │
# │ | Scalability | 8 | 7 | 9 | 7 |        │
# │ | Security | 9 | 8 | 7 | 8 |          │
# │ | Integration | 7 | 9 | 6 | 5 |        │
# │ | Reputation | 8 | 7 | 8 | 6 |        │
# │                                          │
# │ ## Weighted Scores                       │
# │                                          │
# │ | Criteria | Score A | Score B | Score C | Score D |      │
# │ |----------|---------|---------|---------|---------|      │
# │ | Cost | 35 | 40 | 30 | 25 |           │
# │ | Features | 32 | 28 | 36 | 24 |        │
# │ | Usability | 24 | 36 | 28 | 32 |      │
# │ | Support | 21 | 18 | 24 | 15 |         │
# │ | Scalability | 32 | 28 | 36 | 28 |     │
# │ | Security | 45 | 40 | 35 | 40 |        │
# │ | Integration | 21 | 27 | 18 | 15 |     │
# │ | Reputation | 24 | 21 | 24 | 18 |     │
# │ |----------|---------|---------|---------|---------|      │
# │ | **TOTAL** | **234** | **246** | **231** | **197** |  │
# │                                          │
# │ ## Recommendation                      │
# │                                          │
# │ **Winner: Vendor B** (246 points)       │
# │                                          │
# │ ### Sensitivity Analysis                │
# │ - If cost weight increases → Vendor A wins │
# │ - If features weight increases → Vendor C wins │
# │                                          │
# │ ### Next Steps                          │
# │ 1. Request demos from top 2 vendors    │
# │ 2. Conduct reference checks             │
# │ 3. POC with selected vendor            │
# │ 4. Negotiate contract                  │
# └─────────────────────────────────────────┘
```

#### Expected Outcomes

```
✓ Quantitative comparison
✓ Weighted scoring system
✓ Sensitivity analysis
✓ Clear recommendation
```

---

## 6. Problem Solving Scenarios

### 6.1 Scenario: Root Cause Analysis

#### Description

User systematically analyzes problems using structured thinking frameworks.

#### User Story

```
As a problem solver
I want to systematically analyze problems
So that I can identify root causes efficiently
```

#### Workflow

```bash
# User has a problem to analyze
$ otr "Analyze why our user engagement is dropping"

# System shows analysis options:
# ┌─────────────────────────────────────────┐
# │ Thinking Chain                          │
# │ [analysis] ──► [product] ──► [engagement] ◉│
# │ [analysis] ──► [product] ──► [retention] │
# │ [analysis] ──► [business] ──► [revenue] │
# │ [analysis] ──► [team] ──► [productivity] │
# │                                         │
# │ Context: User engagement metric         │
# └─────────────────────────────────────────┘

# User selects "engagement" branch
# [Space] Generate

# Output:
# ┌─────────────────────────────────────────┐
# │ Generated Analysis Framework              │
# │ =====================================  │
# │                                          │
# │ # Root Cause Analysis: User Engagement │
# │                                          │
# │ ## 1. Problem Statement               │
# │                                          │
# │ **Issue**: User engagement dropped 30% │
# │ **Timeline**: Last 30 days             │
# │ **Impact**: Revenue, retention         │
# │ **Severity**: High                     │
# │                                          │
# │ ## 2. Data Evidence                   │
# │                                          │
# │ ### Quantitative Data                  │
# │ - DAU dropped from 10K to 7K          │
# │ - Session time: 5min → 3min           │
# │ - Feature usage: -40%                  │
# │ - Return rate: 20% → 35%              │
# │                                          │
# │ ### Qualitative Data                   │
# │ - Support tickets: +50%                │
# │ - NPS score: 45 → 32                  │
# │ - App reviews: 2.8★ → 2.3★         │
# │                                          │
# │ ## 3. Hypothesis Generation           │
# │                                          │
# │ ### Hypothesis 1: Recent UI Change    │
# │ - Evidence: Engagement dropped after v2.5 │
# │ - Test: Compare cohorts before/after  │
# │ - Data needed: Heatmap, user feedback │
# │                                          │
# │ ### Hypothesis 2: Competitor Launch    │
# │ - Evidence: Competitor X launched       │
# │ - Test: Traffic source analysis       │
# │ - Data needed: Referral data           │
# │                                          │
# │ ### Hypothesis 3: Seasonal Pattern    │
# │ - Evidence: Q1 typically slower        │
# │ - Test: YoY comparison                │
# │ - Data needed: Historical data         │
# │                                          │
# │ ## 4. Root Cause Verification        │
# │                                          │
# │ ### Priority 1: UI Change Impact      │
# │ - Investigation: Heatmap analysis     │
# │ - Findings: [Pending data]            │
# │ - Conclusion: [Pending]               │
# │                                          │
# │ ### Priority 2: Competitor Impact     │
# │ - Investigation: Source analysis       │
# │ - Findings: [Pending data]            │
# │ - Conclusion: [Pending]               │
# │                                          │
# │ ## 5. Recommendations               │
# │                                          │
# │ ### Immediate Actions                  │
# │ - [ ] Roll back UI for test group    │
# │ - [ ] Increase user research           │
# │ - [ ] Monitor key metrics daily       │
# │                                          │
# │ ### Short-term Fixes                  │
# │ - [ ] Fix reported friction points    │
# │ - [ ] A/B test key changes            │
# │                                          │
# │ ### Long-term Improvements            │
# │ - [ ] Better change management       │
# │ - [ ] More user research before launch│
# │ - [ ] Improved metrics monitoring    │
# └─────────────────────────────────────────┘

# User actions:
# [Enter] Copy to clipboard
# [y] Save reflex "engagement-analysis"
```

#### Expected Outcomes

```
✓ Structured problem analysis
✓ Data-driven hypothesis generation
✓ Clear verification steps
✓ Actionable recommendations
```

### 6.2 Scenario: Troubleshooting Framework

#### Description

User troubleshoots issues using systematic checklists.

#### User Story

```
As a troubleshooter
I want a systematic approachSo that I can```

#### Workflow

 to problem-solving
 diagnose issues efficiently
```bash
# User has an issue to troubleshoot
$ otr "Website is running slow"

# System shows troubleshooting options:
# ┌─────────────────────────────────────────┐
# │ Thinking Chain                          │
# │ [slow] ──► [website] ──► [frontend] ◉│
# │ [slow] ──► [website] ──► [backend] │
# │ [slow] ──► [website] ──► [database] │
# │ [slow] ──► [website] ──► [network] │
# │                                         │
# │ Context: User-reported slow response   │
# └─────────────────────────────────────────┘

# User selects "frontend" branch
# [Space] Generate

# Output:
# ┌─────────────────────────────────────────┐
# │ Generated Troubleshooting Checklist       │
# │ =====================================  │
# │                                          │
# │ # Slow Website Troubleshooting         │
# │                                          │
# │ ## 1. Quick Checks (First 5 min)    │
# │                                          │
# │ ### Performance Baseline              │
# │ □ Check current load time: [Tool]    │
# │ □ Compare to baseline: [Historical]  │
# │ □ Identify affected pages: [List]     │
# │ □ Check time of day pattern: [Data]  │
# │                                          │
# │ ### External Factors                  │
# │ □ Check CDN status: [Dashboard]      │
# │ □ Check DNS resolution: [Tool]       │
# │ □ Check SSL certificate: [Status]    │
# │ □ Verify internet connection: [Local]│
# │                                          │
# │ ## 2. Frontend Analysis             │
# │                                          │
# │ ### Page Load Metrics                 │
# │ □ First Contentful Paint (FCP): [Target: <1.8s]│
# │ □ Largest Contentful Paint (LCP): [Target: <2.5s]│
# │ □ Time to Interactive (TTI): [Target: <3.8s]│
# │ □ Total Blocking Time (TBT): [Target: <200ms]│
# │                                          │
# │ ### Resource Analysis                 │
# │ □ Check image sizes: [DevTools]      │
# │ □ Check JS bundle size: [Lighthouse]│
# │ □ Check CSS complexity: [Coverage]   │
# │ □ Check third-party scripts: [Tags] │
# │                                          │
# │ ## 3. Code Issues                    │
# │                                          │
# │ ### JavaScript                        │
# │ □ Heavy computations: [Profiler]    │
# │ □ Memory leaks: [Memory tab]         │
# │ □ Inefficient loops: [Analysis]       │
# │ □ Blocking operations: [Async check] │
# │                                          │
# │ ### Network Requests                  │
# │ □ Too many requests: [Network tab]  │
# │ □ Large payloads: [Size analysis]   │
# │ □ Slow third-party: [Timing]         │
# │ □ Unnecessary requests: [Coverage]   │
# │                                          │
# │ ## 4. Recommended Actions            │
# │                                          │
# │ ### Quick Wins                        │
# │ 1. Compress images                    │
# │ 2. Enable compression                  │
# │ 3. Remove unused code                 │
# │ 4. Lazy load below-fold content      │
# │                                          │
# │ ### Medium-term Fixes                 │
# │ 1. Code split bundles                 │
# │ 2. Optimize third-party               │
# │ 3. Implement caching                  │
# │                                          │
# │ ### Long-term Improvements            │
# │ 1. Migration to faster framework    │
# │ 2. Server-side rendering             │
# │ 3. Edge caching                      │
# └─────────────────────────────────────────┘

# User actions:
# [Enter] Copy to clipboard
# [c] Create investigation ticket
```

#### Expected Outcomes

```
✓ Systematic approach
✓ Time-boxed checks
✓ Quantifiable metrics
✓ Prioritized actions
```

---

## 7. Learning & Research Scenarios

### 7.1 Scenario: Concept Explanation

#### Description

Learner gets structured explanations of new concepts.

#### User Story

```
As a learner
I want clear explanations of new concepts
So that I can understand quickly
```

#### Workflow

```bash
# User wants to learn something
$ otr "Explain machine learning to me"

# Output:
# ┌─────────────────────────────────────────┐
# │ AI Response                              │
# │ =====================================  │
# │                                          │
# │ # Machine Learning: A Beginner's Guide │
# │                                          │
# │ ## What Is Machine Learning?           │
# │                                          │
# │ Machine Learning (ML) is a subset of   │
# │ artificial intelligence that enables   │
# │ systems to learn and improve from    │
# │ experience without being explicitly   │
# │ programmed.                            │
# │                                          │
# │ ## Key Concepts                       │
# │                                          │
# │ ### 1. Training Data                  │
# │ - Data used to teach the model       │
# │ - Examples with known answers         │
# │ - Quality affects model performance   │
# │                                          │
# │ ### 2. Model                          │
# │ - Mathematical representation of patterns│
# │ - Learns from training data           │
# │ - Makes predictions on new data       │
# │                                          │
# │ ### 3. Features                       │
# │ - Input variables used for prediction │
# │ - Selected based on domain expertise │
# │ - Can be engineered from raw data    │
# │                                          │
# │ ### 4. Labels                        │
# │ - Expected output/target variable      │
# │ - Known answers in training data     │
# │ - What we want to predict            │
# │                                          │
# │ ## Types of Machine Learning        │
# │                                          │
# │ | Type | Data | Supervision | Use Case |   │
# │ |------|-------|-------------|----------|   │
# │ | Supervised | Labeled | Full | Prediction | │
# │ | Unsupervised | Unlabeled | None | Patterns | │
# │ | Reinforcement | Rewards | Partial | Decisions | │
# │                                          │
# │ ## Simple Example                      │
# │                                          │
# │ # Predicting house prices
```
 │
# │ # Training data:
# │ # [size, rooms, price]
# │ data = [
# │   ([100, 2], 200000),
# │   ([150, 3], 300000),
# │   ([200, 4], 400000),
# │ ]
# │
# │ # Model learns: bigger = pricier
# │ # New prediction:
# │ predict([175, 3]) → ~350000
# │ ```
# │
# │ ## Real-World Applications
# │
# │ - Email spam detection
# │ - Product recommendations
# │ - Medical diagnosis
# │ - Fraud detection
# │ - Voice assistants
# │
# │ ## Learning Path
# │
# │ 1. **Statistics**: Understand basics
# │ 2. **Python**: Learn the language
# │ 3. **Scikit-learn**: Start simple
# │ 4. **Deep Learning**: Neural nets
# │ 5. **Specialize**: Pick your domain
# │
# │ ## Common Misconceptions
# │
# │ ❌ ML = AI (It's a subset)
# │ ❌ More data = Always better
# │ ❌ ML can solve any problem
# │ ❌ Models are always right
# │
# │ ## Next Steps
# │
# │ - Try: Google Colab tutorials
# │ - Read: "Hands-On ML" book
# │ - Practice: Kaggle competitions
# └─────────────────────────────────────────┘

# User actions:
# [Enter] Copy to clipboard
# [y] Save reflex "ml-basics"
```

#### Expected Outcomes

```
✓ Clear explanations
✓ Progressive complexity
✓ Examples and analogies
✓ Learning path provided
```

### 7.2 Scenario: Comparative Study

#### Description

Learner compares technologies, methods, or concepts.

#### User Story

```
As a learner
I want to compare different options
So that I can make informed choices
```

#### Workflow

```bash
# User wants to compare options
$ otr "Compare React vs Vue vs Angular"

# Output:
# ┌─────────────────────────────────────────┐
# │ Comparative Analysis                      │
# │ =====================================  │
# │                                          │
# │ # React vs Vue vs Angular: Comparison  │
# │                                          │
# │ ## Overview                             │
# │                                          │
# │ | Framework | Created | By | Philosophy |   │
# │ |-----------|---------|----|------------|   │
# │ | React | 2013 | Meta (Facebook) | UI as functions |   │
# │ | Vue | 2014 | Evan You | Progressive |   │
# │ | Angular | 2010 | Google | Full framework |   │
# │                                          │
# │ ## Learning Curve                      │
# │                                          │
# │ | Factor | React | Vue | Angular |       │
# │ |--------|-------|-----|--------|       │
# │ | Difficulty | Medium | Easy | Hard |   │
# │ | Setup | Fast | Fast | Slow |         │
# │ | Concepts | 1-2 core | 2-3 core | Many |   │
# │ | TS Support | Good | Excellent | Native |   │
# │                                          │
# │ ## Performance                          │
# │                                          │
# │ | Metric | React | Vue | Angular |      │
# │ |--------|-------|-----|--------|      │
# │ | Bundle Size | Small | Smallest | Large |   │
# │ | Initial Load | Fast | Fast | Slow |   │
# │ | Update Speed | Fast | Faster | Fast |   │
# │ | Memory | Low | Lower | Higher |      │
# │                                          │
# │ ## Developer Experience                │
# │                                          │
# │ | Factor | React | Vue | Angular |     │
# │ |--------|-------|-----|--------|     │
# │ | Dev Tools | Excellent | Excellent | Good | │
# │ | Community | Largest | Large | Large |   │
# │ | Job Market | Highest | Growing | Large | │
# │ | Flexibility | Highest | High | Low |  │
# │                                          │
# │ ## Best For                            │
# │                                          │
# │ **React**:                             │
# │ - Large applications                   │
# │ - Teams needing flexibility           │
# │ - Complex state management            │
# │ - Existing React expertise            │
# │                                          │
# │ **Vue**:                              │
# │ - Quick prototyping                   │
# │ - Teams new to frameworks             │
# │ - When you want best DX out of box   │
# │ - Clean, readable code preference     │
# │                                          │
# │ **Angular**:                          │
# │ - Enterprise applications             │
# │ - Teams preferring structure          │
# │ - When you need everything included   │
# │ - Long-term projects                  │
# │                                          │
# │ ## Decision Factors                   │
# │                                          │
# │ ### Choose React if:                  │
# │ □ Your team knows React               │
# │ □ Flexibility is important             │
# │ □ Complex interactions needed         │
# │ □ Large ecosystem required            │
# │                                          │
# │ ### Choose Vue if:                    │
# │ □ Quick to learn                      │
# │ □ Clean code is priority             │
# │ □ Incremental adoption               │
# │ □ TypeScript is important            │
# │                                          │
# │ ### Choose Angular if:                │
# │ □ Enterprise scale                    │
# │ □ Structure is valued                │
# │ □ Full framework needed              │
# │ □ Long-term commitment               │
# │                                          │
# │ ## 2024 Trends                        │
# │                                          │
# │ - React: Server Components, AI integration
# │ - Vue: Vapor Mode, improved performance
# │ - Angular: Signals, stricter typing
# │                                          │
# │ ## Recommendation                      │
# │                                          │
# │ For most new projects: **Vue** for    │
# │ beginners, **React** for flexibility.  │
# │                                          │
# │ For enterprise: **Angular** if structure
# │ is valued.                             │
# └─────────────────────────────────────────┘

# User actions:
# [Enter] Copy to clipboard
# [y] Save reflex "frontend-compare"
```

#### Expected Outcomes

```
✓ Structured comparison
✓ Objective criteria
✓ Clear recommendations
✓ Updated for current year
```

---

## 8. Documentation Scenarios

### 8.1 Scenario: Project Documentation Generation

#### Description

Writer generates comprehensive project documentation.

#### User Story

```
As a technical writer
I want to generate project documentation
So that I can maintain consistent documentation structure
```

#### Workflow

```bash
# User triggers documentation generation
$ otr "Create project documentation"

# System shows:
# ┌─────────────────────────────────────────┐
# │ Thinking Chain                          │
# │ [docs] ──► [project] ──► [technical] ◉│
# │ [docs] ──► [project] ──► [business] │
# │ [docs] ──► [process] ──► [onboarding]│
# │ [docs] ──► [process] ──► [runbook] │
# │                                         │
# │ Confidence: 78%                          │
# └─────────────────────────────────────────┘

# User selects "technical" branch
# [Space] Generate

# Output:
# ┌─────────────────────────────────────────┐
# │ Generated Documentation                  │
# │ =====================================  │
# │                                          │
# │ # Project Documentation                │
# │                                          │
# │ ## 1. Overview                         │
# │    - What is this project?             │
# │    - Why does it exist?                │
# │    - Target audience                   │
# │                                          │
# │ ## 2. Quick Start                     │
# │    - Prerequisites                     │
# │    - Installation steps                 │
# │    - Basic usage example               │
# │                                          │
# │ ## 3. Core Concepts                  │
# │    - Key terminology                   │
# │    - Architecture overview             │
# │    - Data flow diagrams                │
# │                                          │
# │ ## 4. Features                        │
# │    - Feature 1: Description          │
# │    - Feature 2: Description          │
# │    - Feature 3: Description          │
# │                                          │
# │ ## 5. Usage Guide                    │
# │    - Step-by-step instructions        │
# │    - Common workflows                 │
# │    - Examples and demos              │
# │                                          │
# │ ## 6. Configuration                  │
# │    - Environment variables            │
# │    - Config file structure            │
# │    - Options reference                │
# │                                          │
# │ ## 7. Troubleshooting                │
# │    - Common issues                   │
# │    - Error messages                  │
# │    - FAQ                             │
# │                                          │
# │ ## 8. Contributing                   │
# │    - Development setup               │
# │    - Coding standards                │
# │    - Pull request process            │
# │                                          │
# │ ## 9. Change Log                     │
# │    - Version history                 │
# │    - Breaking changes                │
# │    - Migration guide                 │
# └─────────────────────────────────────────┘

# [Enter] Save to docs/
# [y] Save reflex "tech-project-docs"
```

#### Expected Outcomes

```
✓ Comprehensive template
✓ Consistent formatting
✓ Easy to customize
✓ Covers common sections
```

### 8.2 Scenario: Process Documentation

#### Description

Writer documents processes and workflows.

#### User Story

```
As a process owner
I want to document our team processes
So that new members can onboard faster
```

#### Workflow

```bash
# User needs process documentation
$ otr "Document our sprint retrospective process"

# System shows:
# ┌─────────────────────────────────────────┐
# │ Thinking Chain                          │
# │ [process] ──► [agile] ──► [ceremony] ◉│
# │ [process] ──► [agile] ──► [planning] │
# │ [process] ──► [devops] ──► [deployment]│
# │ [process] ──► [onboarding] ──► [checklist]│
# │                                         │
# │ Context: Sprint retrospective           │
# └─────────────────────────────────────────┘

# User selects "agile ceremony" branch
# [Space] Generate

# Output:
# ┌─────────────────────────────────────────┐
# │ Generated Process Documentation           │
# │ =====================================  │
# │                                          │
# │ # Sprint Retrospective Process          │
# │                                          │
# │ ## Purpose                             │
# │ The sprint retrospective is a meeting   │
# │ at the end of each sprint where the   │
# │ team reflects on how the sprint went  │
# │ and identifies improvements for the   │
# │ next sprint.                          │
# │                                          │
# │ ## Participants                        │
# │ - Core team members                    │
# │ - Scrum Master (facilitator)          │
# │ - Product Owner (optional)            │
# │                                          │
# │ ## Timing                              │
# │ - Duration: 60 minutes                │
# │ - Frequency: End of each sprint      │
# │ - When: Thursday, 3:00 PM            │
# │                                          │
# │ ## Pre-requisites                     │
# │ □ Sprint completed (all stories done)  │
# │ □ Metrics collected                   │
# │ □ Impediments logged                  │
# │ □ Team availability confirmed          │
# │                                          │
# │ ## Agenda                             │
# │                                          │
# │ ### 1. Opening (5 min)               │
# │ - Welcome and goal setting            │
# │ - Review of previous action items     │
# │                                          │
# │ ### 2. Data Gathering (15 min)        │
# │ - Review sprint metrics               │
# │ - Individual reflections              │
# │ - Highlight achievements              │
# │                                          │
# │ ### 3. Generate Insights (20 min)     │
# │ - What went well?                    │
# │ - What didn't go well?               │
# │ - Surprises and learnings             │
# │                                          │
# │ ### 4. Decide Actions (15 min)        │
# │ - Prioritize improvements            │
# │ - Define action items                 │
# │ - Assign owners and deadlines         │
# │                                          │
# │ ### 5. Closing (5 min)                │
# │ - Summarize key points               │
# │ - Confirm action items                │
# │ - Next sprint preview                │
# │                                          │
# │ ## Facilitation Tips                   │
# │ - Create safe environment             │
# │ - Focus on process, not people       │
# │ - Encourage participation              │
# │ - Keep it positive                   │
# │                                          │
# │ ## Common Pitfalls                    │
# │ - Turning into blame session         │
# │ - Not following up on actions        │
# │ - Taking too long                    │
# │ - Skipping the meeting               │
# │                                          │
# │ ## Action Item Template              │
# │ - What: [Specific action]           │
# │ - Who: [Owner]                      │
# │ - When: [Deadline]                   │
# │ - How: [Success criteria]            │
# │                                          │
# │ ## Metrics to Track                  │
# │ - Number of action items             │
# │ - Completion rate                    │
# │ - Improvement satisfaction           │
# │ - Sprint velocity trend              │
# └─────────────────────────────────────────┘

# [Enter] Save to docs/processes/
# [y] Save reflex "sprint-retro-process"
```

#### Expected Outcomes

```
✓ Clear process flow
✓ Time-boxed agenda
✓ Facilitation tips
✓ Common pitfalls awareness
```

---

## 9. Edge Cases

### 9.1 Scenario: No Matching Reflex

#### Description

User input doesn't match any existing reflex.

#### Workflow

```bash
$ otr "Explain quantum computing"

# Output:
# ┌─────────────────────────────────────────┐
# │ Thinking Chain                          │
# │ [No matching reflex found]              │
# │                                          │
# │ Total reflections: 15                    │
# │ Searched in: default space              │
# └─────────────────────────────────────────┘
# ┌─────────────────────────────────────────┐
# │ Message                                  │
# │ No matching reflex found.                │
# │                                          │
# │ [Space] Generate with AI                 │
# │ [Esc] Cancel                             │
# │                                          │
# │ Tip: Create a reflex after generation    │
# │     for faster access next time          │
# └─────────────────────────────────────────┘

# [Space] Generate

# Output:
# ┌─────────────────────────────────────────┐
# │ AI Response                              │
# │ =====================================  │
# │                                          │
# │ # Quantum Computing Basics              │
# │                                          │
# │ Quantum computing uses quantum bits       │
# │ (qubits) instead of classical bits.      │
# │                                          │
# │ ## Key Concepts                         │
# │                                          │
# │ - **Superposition**: Qubits can be 0     │
# │   and 1 simultaneously                  │
# │ - **Entanglement**: Qubits can be       │
# │   connected (spooky action at a         │
# │   distance)                              │
# │ - **Quantum Gate**: Operations on       │
# │   qubits (like NOT, CNOT, Hadamard)   │
# │                                          │
# │ ## Potential Applications                │
# │                                          │
# │ - Cryptography (breaking RSA)           │
# │ - Drug discovery                         │
# │ - Optimization problems                 │
# │ - Machine learning                     │
# │                                          │
# │ ## Current State                        │
# │                                          │
# │ IBM, Google, and others have quantum     │
# │ computers, but they're still early       │
# │ stage with high error rates.            │
# └─────────────────────────────────────────┘

# [y] Save as reflex "quantum-basics"
```

#### Expected Outcomes

```
✓ Clear "no match" feedback
✓ Immediate AI generation option
✓ Helpful tip about creating reflexes
✓ Optional save as reflex
```

### 9.2 Scenario: Multiple Reflex Matches

#### Description

User input matches multiple reflexes with similar confidence.

#### Workflow

```bash
# User has two similar reflexes
$ otr "analysis"

# Output:
# ┌─────────────────────────────────────────┐
# │ Thinking Chain (Multiple Matches)        │
# ├─────────────────────────────────────────┤
# │ [analysis] ──► [swot] ──► [business] ✓ 85%│
# │ [analysis] ──► [root-cause] ──► [problem] 82%│
# │ [analysis] ──► [comparison] ──► [options] 78%│
# │ [analysis] ──► [metrics] ──► [kpis] 65%│
# │                                          │
# │ Multiple matches found!                  │
# │ Use [↑/↓] to select, [Space] to execute│
# └─────────────────────────────────────────┘

# User selects "comparison" branch
# [↓] [→] Select "comparison" → "options"
# [Space] Execute
```

#### Expected Outcomes

```
✓ All matches displayed
✓ Confidence scores shown
✓ Easy navigation with arrow keys
✓ Clear visual selection state
```

### 9.3 Scenario: Low Confidence Match

#### Description

System finds a match but confidence is too low.

#### Workflow

```bash
$ otr "build cool app quickly"

# Output:
# ┌─────────────────────────────────────────┐
# │ Thinking Chain                          │
# │ [build app] ──► [outline] ⚠️ 45%    │
# │                                          │
# │ ⚠️ Low confidence match (45%)          │
# │                                          │
# │ Suggestions:                            │
# │ 1. Add more context in your query      │
# │ 2. Use more specific keywords           │
# │ 3. Create a new reflex after generation │
# └─────────────────────────────────────────┘
# ┌─────────────────────────────────────────┐
# │ Options                                 │
# │                                          │
# │ [Space] Proceed with this reflex        │
# │ [↓] Show more options                  │
# │ [Esc] Cancel, use more specific query  │
# └─────────────────────────────────────────┘

# User tries more specific query
# [Esc] Cancel
$ otr "Create project proposal outline for new product"

# System finds strong match:
# ┌─────────────────────────────────────────┐
# │ Thinking Chain                          │
# │ [proposal outline] ──► [product] ✓ 92%│
# │                                          │
# │ ✓ High confidence match (92%)          │
# └─────────────────────────────────────────┘
```

#### Expected Outcomes

```
✓ Low confidence warning
✓ Helpful suggestions
✓ Option to cancel and refine
✓ Improved results with better query
```

### 9.4 Scenario: Expired/Decayed Reflex

#### Description

User tries to use a reflex that has decayed due to non-use.

#### Workflow

```bash
# User tries old reflex
$ otr "old-pattern"

# Output:
# ┌─────────────────────────────────────────┐
# │ Thinking Chain                          │
# │ [old-pattern] ──► [v1] ⚠️ DECAYED   │
# │                                          │
# │ ⚠️ This reflex has decayed!           │
# │                                          │
# │ Original Strength: 100%                 │
# │ Current Strength: 15% (below threshold)  │
# │ Last Used: 90 days ago                  │
# │ Decay Count: 5                         │
# │                                          │
# │ [Space] Use anyway (low confidence)    │
# │ [r] Refresh reflex (re-learn)           │
# │ [d] Delete this reflex                 │
# └─────────────────────────────────────────┘

# User chooses to refresh
# [r] Refresh

# Output:
# ┌─────────────────────────────────────────┐
# │ Reflex Refresh                          │
# │ =====================================  │
# │                                          │
# │ Pattern: old-pattern                    │
# │ Current Strength: 15%                   │
# │                                          │
# │ [Space] Generate fresh content          │
# │ The new generation will reset strength  │
# │ to 10% and replace old content         │
# └─────────────────────────────────────────┘

# [Space] Generate
```

#### Expected Outcomes

```
✓ Decay status clearly shown
✓ Options to use, refresh, or delete
✓ Graceful degradation
✓ Opportunity to re-learn
```

---

## 10. Scenario Comparison Matrix

### 10.1 Scenario Characteristics

| Scenario | Reflex Strength | Time Saved | Complexity | User Level |
|----------|-----------------|------------|-------------|------------|
| First-Time Onboarding | N/A (new) | Setup time | Medium | Beginner |
| Rapid Reflex Execution | High (80%+) | 30-60s | Low | Intermediate |
| Strength Progression | Growing | Variable | Low | Beginner |
| Article Outline | Medium (50%+) | 5-10min | Medium | Intermediate |
| SWOT Analysis | Medium | 10-20min | Medium | Intermediate |
| Root Cause Analysis | Medium | 15-30min | High | Advanced |
| Project Documentation | Medium | 10-20min | Low | Beginner |
| Concept Explanation | N/A (ad-hoc) | 3-5min | Low | All |
| Comparative Study | Low | 5-10min | Low | Beginner |

### 10.2 Workflow Patterns

| Pattern | Description | Example Scenarios |
|---------|-------------|-------------------|
| **Direct Execute** | Trigger reflex → One-key execute | Rapid Execution |
| **Select & Execute** | Choose branch → Execute | Analysis, Documentation |
| **Interactive Learning** | Step-by-step progression | Concept Learning |
| **Diagnostic** | Checklist → Resolution | Root Cause Analysis |
| **Template Fill** | Select template → Generate | Documentation |

### 10.3 User Journey Map

```
User Journey Stages
================================================================================

Stage 1: Discovery
─────────────────────────────────────────────────────────────────────────────
│ Activities:                                                               │
│ - Install CLI                                                            │
│ - Run initial setup                                                      │
│ - Ask first question                                                     │
│                                                                          │
│ Touchpoints:                                                             │
│ - otr init                                                               │
│ - First AI generation                                                   │
│                                                                          │
│ Success Criteria:                                                        │
│ ✓ Configuration complete                                                │
│ ✓ First generation successful                                            │
└─────────────────────────────────────────────────────────────────────────────

Stage 2: Learning
─────────────────────────────────────────────────────────────────────────────
│ Activities:                                                               │
│ - Explore capabilities                                                   │
│ - Create first reflex                                                   │
│ - Learn keyboard shortcuts                                               │
│                                                                          │
│ Touchpoints:                                                             │
│ - otr help                                                              │
│ - Pattern creation commands                                             │
│ - Reflex status checks                                                  │
│                                                                          │
│ Success Criteria:                                                        │
│ ✓ 3+ reflexes created                                                   │
│ ✓ Comfortable with shortcuts                                             │
│ ✓ Understands strength system                                            │
└─────────────────────────────────────────────────────────────────────────────

Stage 3: Adoption
─────────────────────────────────────────────────────────────────────────────
│ Activities:                                                              │
│ - Daily use of common reflexes                                           │
│ - Build personal reflex library                                          │
│ - Share patterns (v2.0)                                                 │
│                                                                          │
│ Touchpoints:                                                            │
│ - Rapid execution workflow                                              │
│ - Reflex statistics                                                     │
│ - Export commands                                                      │
│                                                                          │
│ Success Criteria:                                                        │
│ ✓ 10+ high-strength reflexes                                             │
│ ✓ Measurable time savings                                               │
│ ✓ Reduced AI prompting                                                 │
└─────────────────────────────────────────────────────────────────────────────

Stage 4: Mastery
─────────────────────────────────────────────────────────────────────────────
│ Activities:                                                              │
│ - Advanced analysis workflows                                            │
│ - Custom reflex creation                                               │
│ - Team knowledge sharing (v2.0)                                         │
│                                                                          │
│ Touchpoints:                                                            │
│ - Complex diagnostic workflows                                          │
│ - Custom pattern extraction                                             │
│ - Skill export                                                         │
│                                                                          │
│ Success Criteria:                                                        │
│ ✓ 50+ reflexes                                                         │
│ ✓ Custom workflows automated                                             │
│ ✓ Team-wide knowledge base                                              │
└─────────────────────────────────────────────────────────────────────────────
```

### 10.4 Command Usage by Scenario

| Scenario | Key Commands Used | Frequency |
|----------|-------------------|-----------|
| Onboarding | `otr init`, `otr` | Once |
| Daily Execution | `otr <trigger>`, `[Space]` | 10-30x/day |
| Pattern Management | `otr pattern create/show/list` | 1-2x/week |
| Documentation | `otr pattern create` | 2-3x/week |
| Analysis | `otr <analysis-type>`, `[Space]` | 5-10x/week |
| Learning | `otr explain <topic>`, `[→]` | 2-3x/week |
| Maintenance | `otr stats`, `otr backup`, `otr decay` | 1x/week |

---

## Appendix A: Quick Reference

### A.1 Essential Commands

```bash
# Core Workflow
otr "your query"              # Main command
[Space]                       # Execute selected
[↑/↓]                         # Select branch
[→]                           # Expand branch
[←]                           # Go back
[Esc]                         # Cancel

# Pattern Management
otr pattern create            # Create new reflex
otr pattern list              # List all reflexes
otr pattern show <id>        # View reflex details
otr pattern delete <id>      # Delete reflex

# System
otr stats                     # View statistics
otr config show              # View configuration
otr health                    # Health check
otr backup                    # Create backup
```

### A.2 Keyboard Shortcuts Reference

| Key | Action | Context |
|-----|--------|---------|
| Space | Execute/Generate | Anywhere |
| ↑/↓ | Select branch | Thinking Chain |
| → | Expand branch | Thinking Chain |
| ← | Go back | Thinking Chain |
| Enter | Copy to clipboard | Output |
| h | Show help | Global |
| q / Ctrl+C | Quit | Global |
| Esc | Cancel | Anywhere |

---

**Document Version**: v2.0  
**Created**: 2026-02-20  
**Project**: open-think-reflex