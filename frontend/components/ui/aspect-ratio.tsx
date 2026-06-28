'use client'

import * as AspectRatioPrimitive from '@radix-ui/react-aspect-ratio'

// Guarantees space reservation before children load to eliminate Cumulative Layout Shift (CLS).
function AspectRatio({
  ...props
}: React.ComponentProps<typeof AspectRatioPrimitive.Root>) {
  return <AspectRatioPrimitive.Root data-slot="aspect-ratio" {...props} />
}

export { AspectRatio }
