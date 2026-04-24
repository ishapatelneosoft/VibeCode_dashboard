You are extending an existing Next.js frontend application.

---

## FEATURE

Board UI Design (Jira-style)

---

## OBJECTIVE

Build a board UI that visually matches a modern task management system (similar to Jira-style boards), based on the provided reference.

---

## DESIGN PRINCIPLES

* Clean, minimal UI
* Card-based layout
* Clear visual hierarchy
* Consistent spacing and alignment
* Subtle shadows and rounded corners

---

## LAYOUT STRUCTURE

Page Layout:

* Full-width horizontal scrollable board
* Fixed header (optional)
* Columns aligned horizontally

Board:

* Flex container (row)
* Horizontal scrolling enabled

---

## COLUMN DESIGN

* Fixed width columns (approx 280px–320px)
* Light background (slightly different from page)
* Rounded corners
* Column header:

  * Title (bold)
  * Subtle divider
* Vertical stacking of tasks
* Scrollable column content (if overflow)

---

## TASK CARD DESIGN

Each task card must include:

* White background
* Rounded corners
* Soft shadow
* Padding (consistent spacing)

Content:

1. Title (primary text, bold)
2. Metadata row:

   * Assignee avatar (circle)
   * Due date (if present)
3. Secondary info:

   * Created by (subtle text)

Spacing:

* Vertical spacing between cards
* Internal padding (12–16px)

---

## VISUAL BEHAVIOR

* Hover effect on task card (slight elevation)
* Cursor pointer for cards
* Smooth transitions

---

## DATA BINDING RULES

* Render exactly as backend provides
* Do NOT modify ordering
* Handle null fields gracefully:

  * No assignee → show placeholder avatar
  * No due date → hide field

---

## RESPONSIVENESS

* Horizontal scroll for smaller screens
* No column stacking
* Maintain layout integrity

---

## STYLING APPROACH

* Use Tailwind CSS (preferred)
* Avoid inline styles
* Use reusable classes/components

---

## COMPONENT RESPONSIBILITY

Board:

* Layout container only

Column:

* Handles column structure
* Receives tasks as props

TaskCard:

* Pure presentational component
* No API logic

---

## STRICT RULES

* Do NOT introduce custom UI patterns
* Do NOT over-design (keep minimal)
* Do NOT mix layout + logic
* Do NOT hardcode data

---

## OUTPUT

1. Updated Board component
2. Column component
3. TaskCard component (styled)
4. Styling implementation

---

## GOAL

Deliver a clean, Jira-like board UI that matches the reference design and is ready for future features (drag-drop, modal, etc.).
