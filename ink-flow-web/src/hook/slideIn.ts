const offset = 100
const duration = 300

const map = new WeakMap<HTMLElement, Animation>()

const ob = new IntersectionObserver((entries) => {
  for (const en of entries) {
    if (en.isIntersecting) {
      const animation = map.get(en.target as HTMLElement)
      if (animation) {
        animation.play()
        ob.unobserve(en.target as HTMLElement)
      }
    }
  }
})

const belowViewport = (el: HTMLElement) => {
  const rect = el.getBoundingClientRect()
  return rect.top > window.innerHeight
}
export default {
  mounted(el: HTMLElement) {
    if (!belowViewport(el)) {
      return
    }
    const animation = el.animate(
      [
        {
          transform: `translateY(${offset}px)`,
          opacity: 0.5,
        },
        {
          transform: `translateY(0px)`,
          opacity: 1,
        },
      ],
      {
        duration: duration,
        easing: 'ease-out',
        fill: 'forwards',
      },
    )
    animation.pause()

    ob.observe(el)
    map.set(el, animation)
  },
  unmounted(el: HTMLElement) {
    ob.unobserve(el)
  },
}
