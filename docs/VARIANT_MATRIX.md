# Variant Compatibility Matrix

Practical variant assignments for components. Less is more.

## Reduced Variant Set

The original spec had 12 variants. That's too many. Here are the 5 that matter:

| Variant | Use Case |
|---------|----------|
| `View` | Embeddable content, flows in layout |
| `Card` | Bordered container, self-contained |
| `Modal` | Centered overlay with backdrop, blocks interaction |
| `Compact` | Minimal footprint, status bars |
| `Collapsible` | Expandable/collapsible sections |

**Removed variants and why:**
- `Inline` - use `View` instead
- `Window` - use `Card` with title
- `Overlay` - use `Modal` without backdrop
- `Fullscreen` - handled by layout, not component
- `Bar` - use `Compact`
- `Sidebar` - layout concern, not component variant
- `Toast` - Toast is a component, not a variant

---

## Component Variant Matrix

| Component | View | Card | Modal | Compact | Collapsible |
|-----------|:----:|:----:|:-----:|:-------:|:-----------:|
| **Output** |
| streamtext | x | | | | |
| message | x | x | | | |
| thinking | x | | | | x |
| tooluse | x | x | | | x |
| codeblock | x | x | x | | |
| status | | | | x | |
| toast | | x | | | |
| diffview | x | | x | | |
| globresults | x | | | | x |
| grepresults | x | | | | x |
| **Input** |
| chat | x | | | | |
| permission | | x | x | | |
| question | | x | x | | |
| planenter | | | x | | |
| planexit | | x | x | | |
| interrupt | | | x | | |
| confirm | | | x | | |
| filepicker | | | x | | |
| todo | x | | | | x |
| taskmgr | x | | x | | |
| session | | | x | | |
| search | | | x | | |
| **Shared** |
| button | x | | | x | |
| chip | x | | | x | |
| card | | x | | | |

---

## Default Variants

Each component has ONE default variant. Only override when you have a good reason.

| Component | Default | Override When |
|-----------|---------|---------------|
| streamtext | View | Never |
| message | View | Card for standalone display |
| thinking | Collapsible | View when always visible |
| tooluse | Card | Collapsible for long outputs |
| codeblock | View | Modal for full-screen editing |
| status | Compact | Never |
| toast | Card | Never |
| permission | Modal | Card when inline in layout |
| question | Modal | Card when inline in layout |
| todo | Collapsible | View in sidebar layouts |

---

## Implementation Notes

```go
type Variant int

const (
    VariantView Variant = iota
    VariantCard
    VariantModal
    VariantCompact
    VariantCollapsible
)

// Each component defines which variants it supports
func (c *StreamText) SupportedVariants() []Variant {
    return []Variant{VariantView}
}

func (c *Permission) SupportedVariants() []Variant {
    return []Variant{VariantCard, VariantModal}
}
```

Components should panic at init time if given an unsupported variant.
