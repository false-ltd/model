<template>
    <div
        ref="wrapRef"
        class="relative w-full h-full select-none"
        :class="ambient ? 'pointer-events-none' : 'cursor-grab active:cursor-grabbing touch-none'"
        @pointerdown="onPointerDown"
        @pointermove="onPointerMove"
        @pointerup="onPointerUp"
        @pointercancel="onPointerUp"
        @wheel.prevent="onWheel"
    >
        <canvas ref="canvasRef" class="block w-full h-full" />

        <!-- Tooltip (interactive mode) -->
        <div
            v-if="!ambient && hoverPoint"
            class="absolute z-10 pointer-events-none bg-default border border-default rounded-lg shadow-lg px-2.5 py-2 max-w-56"
            :style="{ left: hoverScreen.x + 'px', top: hoverScreen.y + 'px', transform: 'translate(-50%, calc(-100% - 14px))' }"
        >
            <div class="text-xs font-semibold text-default truncate">{{ hoverPoint.n }}</div>
            <div class="text-[10px] text-muted font-mono truncate">{{ providerName(hoverPoint.p) }}</div>
            <div class="flex gap-2 mt-1 text-[10px] text-toned tabular">
                <span v-if="hoverPoint.ci != null">{{ hoverPoint.ci === 0 ? 'FREE' : `$${hoverPoint.ci}/M` }}</span>
                <span v-if="hoverPoint.ctx">{{ formatTokens(hoverPoint.ctx) }} ctx</span>
            </div>
        </div>

        <slot />
    </div>
</template>

<script setup lang="ts">
    /**
     * The Atlas — every indexed model as a star.
     * x: log price (free at left edge, unpriced in the left gutter),
     * y: log context window, color: provider, glow: capability count.
     *
     * ambient=true renders a slow auto-drifting starfield (hero backdrop,
     * no interaction); ambient=false is fully interactive: drag to pan with
     * inertia, wheel/pinch to zoom, click a star to open the model.
     */
    const props = withDefaults(defineProps<{ ambient?: boolean; marginal?: boolean }>(), {
        ambient: false,
        marginal: false,
    });

    const localePath = useLocalePath();
    const lite = useAtlasLite();

    const wrapRef = ref<HTMLElement | null>(null);
    const canvasRef = ref<HTMLCanvasElement | null>(null);

    const { data: atlasData } = await useAtlasData();

    interface Pt {
        id: number; n: string; p: string;
        ci: number | null; co: number | null; ctx: number | null;
        k: number; ow: boolean;
        x: number; y: number; z: number; // normalized world coords + glow weight
        d: number; // world depth for the 3D view
    }

    const providerName = (id: string) =>
        atlasData.value?.data?.providers?.find((p: any) => p.id === id)?.name || id;

    // ---------- world layout ----------
    const points = ref<Pt[]>([]);
    const projBuf: { p: Pt; x: number; y: number; scale: number; depth: number }[] = [];

    const buildPoints = () => {
        const raw: any[] = atlasData.value?.data?.points || [];
        const providers: any[] = atlasData.value?.data?.providers || [];
        const out: Pt[] = raw.map((m) => {
            const price = m.ci ?? null;
            const x = price == null ? 0.02 + hash01(m.id) * 0.03 : 0.06 + Math.log10(1 + price) / Math.log10(1 + 750) * 0.92;
            const ctx = m.ctx ?? null;
            const y = ctx == null ? 0.05 + hash01(m.id + 7) * 0.04 : 0.12 + Math.log10(1024 + ctx) / Math.log10(1024 + 2_100_000) * 0.86;
            return {
                id: m.id, n: m.n, p: m.p, ci: m.ci, co: m.co, ctx: m.ctx, k: m.k ?? 0, ow: !!m.ow,
                x, y: 1 - y, // canvas y grows downward
                z: 0.35 + ((m.k ?? 0) / 6) * 0.65,
                // Depth axis for the 3D view: capability bands in depth, plus
                // a deterministic per-model spread so same-spec models
                // separate instead of overplotting.
                d: (m.k ?? 0) / 6 * 0.7 - 0.35 + (hash01(m.id * 3 + 11) - 0.5) * 0.3,
            };
        });
        points.value = out;
        buildMarginals();
    };

    // ---------- palette ----------
    const colorMode = useColorMode();
    // Palette is cached — getComputedStyle forces style recalc and must
    // never run per frame. Invalidated on theme flips.
    let cachedPalette: { dark: boolean; base: string[] } | null = null;
    const palette = () => {
        if (cachedPalette) return cachedPalette;
        const dark = colorMode.value === "dark";
        const base = [0, 1, 2, 3, 4, 5].map((i) => getCSSVar(`--chart-${i}`));
        cachedPalette = { dark, base };
        return cachedPalette;
    };
    // Re-render on theme flips (colors come from CSS variables).
    watch(
        () => colorMode.value,
        () => {
            sprites.clear();
            cachedPalette = null;
            needsRender = true;
        },
    );

    // Pre-rendered glow sprites per (colorIndex × radius tier) — drawImage is
    // far cheaper than 4k shadowBlur arcs per frame.
    const sprites = new Map<string, HTMLCanvasElement>();
    const getSprite = (color: string, r: number, ow = false) => {
        const key = `${color}|${r}|${ow ? 1 : 0}`;
        let s = sprites.get(key);
        if (s) return s;
        // Radial-gradient glow baked per sprite: ~half the fill area of a
        // stacked-disc halo at a softer peak. Lite mode shrinks it further
        // (3.2r, weaker alpha) for low-spec machines; the open-weights ring
        // always gets the canvas size it needs.
        const size = Math.ceil(Math.max((lite.value ? 3.2 : 4.2) * r, ow ? (r * 1.8 + 1) * 2 : 0));
        s = document.createElement("canvas");
        s.width = s.height = size;
        const c = s.getContext("2d")!;
        const cx = size / 2;
        const glow = c.createRadialGradient(cx, cx, r * 0.5, cx, cx, size / 2);
        glow.addColorStop(0, color);
        // providerColor() returns "hsl(h, s%, l%)" — the alpha variant uses
        // the comma form, which canvas accepts as-is.
        glow.addColorStop(1, color.replace(")", ", 0)"));
        c.globalAlpha = lite.value ? 0.22 : 0.3;
        c.fillStyle = glow;
        c.fillRect(0, 0, size, size);
        c.globalAlpha = 1;
        c.fillStyle = color;
        c.beginPath();
        c.arc(cx, cx, r * 0.58, 0, Math.PI * 2);
        c.fill();
        if (ow) {
            // Open-weights ring baked into the sprite — stroking 3k+ rings
            // per frame was the single largest render cost.
            c.strokeStyle = color;
            c.globalAlpha = lite.value ? 0.5 : 0.65;
            c.lineWidth = 1;
            c.beginPath();
            c.arc(cx, cx, r * 1.8, 0, Math.PI * 2);
            c.stroke();
            c.globalAlpha = 1;
        }
        sprites.set(key, s);
        return s;
    };

    // ---------- camera ----------
    const cam = { x: 0.5, y: 0.5, zoom: 1 };
    const vel = { x: 0, y: 0 };
    let width = 0, height = 0, dpr = 1;
    let raf = 0;
    let running = false;
    let needsRender = true;
    const reduced = () => window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    let ambientT = 0;

    const M_TOP = () => (props.marginal ? 40 : 0);
    const M_LEFT = () => (props.marginal ? 38 : 0);

    // 3D galaxy view: tilt the plot plane around the X axis and apply a
    // gentle perspective. Depth = capability band + per-model spread, so
    // same-spec models separate instead of stacking.
    const is3D = ref(true);
    const TILT_DEFAULT = 0.5;
    const tilt = ref(TILT_DEFAULT);

    const project = (p: Pt) => {
        const iw = Math.max(1, width - M_LEFT());
        const ih = Math.max(1, height - M_TOP());
        const wx = (p.x - cam.x) * cam.zoom + 0.5;
        const wy = (p.y - cam.y) * cam.zoom + 0.5;
        if (!is3D.value) {
            return { x: M_LEFT() + wx * iw, y: M_TOP() + wy * ih, scale: 1, depth: 0 };
        }
        const cosT = Math.cos(tilt.value);
        const sinT = Math.sin(tilt.value);
        const yRot = wy * cosT - p.d * sinT;
        const zRot = wy * sinT + p.d * cosT;
        const focal = 3.2;
        const scale = focal / (focal + zRot * 0.9);
        return {
            x: M_LEFT() + wx * scale * iw,
            y: M_TOP() + yRot * scale * ih,
            scale,
            depth: zRot,
        };
    };

    // ---------- marginal histograms (price along top, context along left) ----------
    const PRICE_BINS = 30;
    const CTX_BINS = 22;
    let priceHist: number[] = [];
    let ctxHist: number[] = [];

    const buildMarginals = () => {
        priceHist = new Array(PRICE_BINS).fill(0);
        ctxHist = new Array(CTX_BINS).fill(0);
        for (const p of points.value) {
            priceHist[Math.min(PRICE_BINS - 1, Math.max(0, Math.floor(p.x * PRICE_BINS)))]++;
            ctxHist[Math.min(CTX_BINS - 1, Math.max(0, Math.floor(p.y * CTX_BINS)))]++;
        }
    };

    const drawMarginals = (ctx: CanvasRenderingContext2D) => {
        if (!props.marginal) return;
        const iw = width - M_LEFT();
        const ih = Math.max(1, height - M_TOP());
        // --ui-primary can be oklch(), so alpha is applied via globalAlpha
        // instead of string parsing.
        const primary = getCSSVar("--ui-primary", "#d97706") || "#d97706";
        const maxP = Math.max(...priceHist, 1);
        const maxC = Math.max(...ctxHist, 1);

        ctx.save();
        ctx.globalAlpha = 0.35;
        ctx.fillStyle = primary;
        // price marginal: bars hang from the top edge of the plot upward
        const bandP = M_TOP() - 8;
        const bwP = iw / PRICE_BINS;
        for (let i = 0; i < PRICE_BINS; i++) {
            const h = (priceHist[i] / maxP) * bandP;
            if (h < 0.5) continue;
            ctx.fillRect(M_LEFT() + i * bwP + 0.5, M_TOP() - h, Math.max(1, bwP - 1), h);
        }
        // context marginal: bars extend leftward from the plot's left edge
        const bandC = M_LEFT() - 8;
        const bhC = ih / CTX_BINS;
        for (let i = 0; i < CTX_BINS; i++) {
            const w = (ctxHist[i] / maxC) * bandC;
            if (w < 0.5) continue;
            ctx.fillRect(M_LEFT() - w, M_TOP() + i * bhC + 0.5, w, Math.max(1, bhC - 1));
        }
        ctx.restore();
        // plot frame hairline
        ctx.save();
        ctx.globalAlpha = 0.3;
        ctx.strokeStyle = primary;
        ctx.lineWidth = 1;
        ctx.strokeRect(M_LEFT() + 0.5, M_TOP() + 0.5, iw - 1, ih - 1);
        ctx.restore();
    };

    buildPoints();
    watch(atlasData, () => {
        buildPoints();
        if (!props.ambient) {
            try {
                render();
            } catch (err) {
                console.error("[atlas] interactive data render failed:", err);
            }
        }
    });

    const render = () => {
        const canvas = canvasRef.value;
        if (!canvas) return;
        const ctx = canvas.getContext("2d");
        if (!ctx) return;
        const { dark, base } = palette();
        ctx.clearRect(0, 0, width, height);
        const iw = Math.max(1, width - M_LEFT());
        const ih = Math.max(1, height - M_TOP());
        let lastAlpha = -1;

        // faint grid rings at decade boundaries
        ctx.strokeStyle = dark ? "rgba(255,255,255,0.05)" : "rgba(0,0,0,0.05)";
        ctx.lineWidth = 1;
        for (let gx = 0; gx <= 4; gx++) {
            const sx = M_LEFT() + ((gx / 4 - cam.x) * cam.zoom + 0.5) * iw;
            ctx.beginPath();
            ctx.moveTo(sx, M_TOP());
            ctx.lineTo(sx, height);
            ctx.stroke();
        }
        for (let gy = 0; gy <= 4; gy++) {
            const sy = M_TOP() + ((gy / 4 - cam.y) * cam.zoom + 0.5) * ih;
            ctx.beginPath();
            ctx.moveTo(M_LEFT(), sy);
            ctx.lineTo(width, sy);
            ctx.stroke();
        }

        const baseR = props.ambient ? 1.1 : 1.35;
        const dim = props.ambient ? 0.55 : 1;

        // Project once, then paint far-to-near so nearer stars overlay.
        projBuf.length = 0;
        for (const p of points.value) {
            const pr = project(p);
            if (pr.x < -30 || pr.x > width + 30 || pr.y < -30 || pr.y > height + 30) continue;
            projBuf.push({ p, x: pr.x, y: pr.y, scale: pr.scale, depth: pr.depth });
        }
        projBuf.sort((a, b) => b.depth - a.depth);
        const hoverP = hoverPoint.value;
        const focusProvider = highlightProvider.value;
        const showLabels = !props.ambient && cam.zoom > 2.4;
        let labelBudget = cam.zoom > 8 ? 90 : 36;

        // Additive blending in dark mode: overlapping stars add up to a
        // brighter glow, so density reads naturally (like a real starfield).
        // Light mode keeps normal blending — additive would white out.
        if (dark) ctx.globalCompositeOperation = "lighter";

        for (const item of projBuf) {
            const p = item.p;
            const s = item;

            const color = providerColor(p.p, dark);
            const r = (baseR + p.z * 2.0) * Math.min(cam.zoom, 4) * item.scale;
            const emphasized = focusProvider ? p.p === focusProvider : false;
            const alpha = focusProvider ? (emphasized ? 1 : 0.12) : hoverP && p.p === hoverP.p ? 1 : dim;

            // Sprite radius is quantized to 0.5px tiers: in 3D the
            // perspective scale is continuous, and keying sprites on raw
            // radii would bake a canvas per model and re-bake on every
            // zoom step. Tiers keep the cache small and stable.
            const rTier = Math.max(1, Math.round(r * 2) / 2);
            const sprite = getSprite(color, rTier, p.ow);
            const size = sprite.width;
            const pointAlpha = alpha * (props.ambient ? 0.8 : 0.92);
            if (pointAlpha !== lastAlpha) {
                ctx.globalAlpha = pointAlpha;
                lastAlpha = pointAlpha;
            }
            ctx.drawImage(sprite, s.x - size / 2, s.y - size / 2);

            if (showLabels && labelBudget > 0 && p.k >= 5) {
                labelBudget--;
                ctx.font = "600 10px Inter, sans-serif";
                ctx.fillStyle = dark ? "rgba(255,255,255,0.75)" : "rgba(0,0,0,0.65)";
                ctx.fillText(p.n, s.x + r + 5, s.y + 3);
            }
        }

        ctx.globalAlpha = 1;
        ctx.globalCompositeOperation = "source-over";
        if (!props.ambient) drawAxes(ctx, dark);
        drawMarginals(ctx);
    };

    const drawAxes = (ctx: CanvasRenderingContext2D, dark: boolean) => {
        ctx.font = "500 10px Inter, sans-serif";
        ctx.fillStyle = dark ? "rgba(255,255,255,0.35)" : "rgba(0,0,0,0.35)";
        const iw = Math.max(1, width - M_LEFT());
        const labels = ["$0", "$1", "$10", "$100", "$750+"];
        for (let i = 0; i < 5; i++) {
            const sx = M_LEFT() + ((i / 4 - cam.x) * cam.zoom + 0.5) * iw;
            if (sx > M_LEFT() + 24 && sx < width - 8) ctx.fillText(labels[i], sx - 8, height - 8);
        }
        if (!props.marginal) {
            ctx.save();
            ctx.translate(12, 16);
            ctx.fillText("1M ctx ↑", 0, 0);
            ctx.restore();
        }
    };

    // ---------- frame loop ----------
    let inertiaRaf = 0;

    // Interactive mode: inertia after a drag, rendered per tick. Pan is
    // kept 1:1 immediate (low latency); only zoom is smoothed.
    const inertiaTick = () => {
        if (Math.abs(vel.x) > 0.00002 || Math.abs(vel.y) > 0.00002) {
            cam.x -= vel.x / cam.zoom;
            cam.y -= vel.y / cam.zoom;
            targetCam.x = cam.x;
            targetCam.y = cam.y;
            clampCam();
            vel.x *= 0.92;
            vel.y *= 0.92;
            render();
            inertiaRaf = requestAnimationFrame(inertiaTick);
        }
    };

    const kickInertia = () => {
        cancelAnimationFrame(inertiaRaf);
        inertiaRaf = requestAnimationFrame(inertiaTick);
    };

    // Ambient mode: drifting starfield loop, paused offscreen.
    const frame = () => {
        if (!running) return;
        if (!reduced()) {
            ambientT += 0.0016;
            cam.x = 0.5 + Math.sin(ambientT * 0.6) * 0.08;
            cam.y = 0.5 + Math.cos(ambientT * 0.42) * 0.06;
        }
        try {
            render();
        } catch (err) {
            running = false;
            console.error("[atlas] ambient render failed:", err);
            return;
        }
        raf = requestAnimationFrame(frame);
    };

    const start = () => {
        if (running) return;
        running = true;
        raf = requestAnimationFrame(frame);
    };
    const stop = () => {
        running = false;
        cancelAnimationFrame(raf);
    };

    // ---------- pointer interaction (interactive mode) ----------
    const hoverPoint = ref<Pt | null>(null);
    const hoverScreen = ref({ x: 0, y: 0 });
    const highlightProvider = ref<string | null>(null);

    let dragging = false;
    let moved = false;
    let lastPtr = { x: 0, y: 0 };
    const pointers = new Map<number, { x: number; y: number }>();
    let pinchDist = 0;

    // Wheel/move bursts coalesce into one repaint per animation frame —
    // synchronous renders per event were the zoom jank.
    let renderQueued = false;
    const scheduleRender = () => {
        if (renderQueued) return;
        renderQueued = true;
        requestAnimationFrame(() => {
            renderQueued = false;
            render();
        });
    };

    // Two-camera model: input mutates targetCam instantly; the rendered
    // cam eases toward it every frame (exponential smoothing, ~90ms). This
    // is what makes zoom feel fluid instead of stepped.
    const targetCam = { x: 0.5, y: 0.5, zoom: 1 };

    const clampCamObj = (c: { x: number; y: number; zoom: number }) => {
        c.zoom = Math.min(20, Math.max(0.8, c.zoom));
        const half = 0.5 / c.zoom;
        c.x = Math.min(1 + half, Math.max(-half, c.x));
        c.y = Math.min(1 + half, Math.max(-half, c.y));
    };
    const clampCam = () => {
        clampCamObj(cam);
        clampCamObj(targetCam);
    };

    let smoothRaf = 0;
    let smoothLast = 0;
    let smoothGuard = 0;
    const snapToTarget = () => {
        cam.x = targetCam.x;
        cam.y = targetCam.y;
        cam.zoom = targetCam.zoom;
        smoothRaf = 0;
    };
    const kickSmooth = () => {
        if (props.ambient || smoothRaf) return;
        smoothLast = performance.now();
        smoothGuard = 0;
        const step = (now: number) => {
            const dt = Math.min(64, now - smoothLast);
            smoothLast = now;
            const k = 1 - Math.exp(-dt / 90);
            cam.x += (targetCam.x - cam.x) * k;
            cam.y += (targetCam.y - cam.y) * k;
            cam.zoom += (targetCam.zoom - cam.zoom) * k;
            if (!isFinite(cam.x) || !isFinite(cam.y) || !isFinite(cam.zoom) || ++smoothGuard > 240) {
                snapToTarget();
                render();
                return;
            }
            render();
            const settled =
                Math.abs(targetCam.zoom - cam.zoom) < 1e-4 &&
                Math.abs(targetCam.x - cam.x) < 1e-6 &&
                Math.abs(targetCam.y - cam.y) < 1e-6;
            if (settled) {
                cam.x = targetCam.x;
                cam.y = targetCam.y;
                cam.zoom = targetCam.zoom;
                smoothRaf = 0;
                render();
                return;
            }
            smoothRaf = requestAnimationFrame(step);
        };
        if (reduced()) {
            cam.x = targetCam.x;
            cam.y = targetCam.y;
            cam.zoom = targetCam.zoom;
            smoothRaf = 0;
            render();
            return;
        }
        smoothRaf = requestAnimationFrame(step);
    };

    const zoomAt = (sx: number, sy: number, factor: number) => {
        // Anchor math compounds on the target so rapid wheel ticks chain
        // cleanly while the camera is still catching up.
        const wx = (sx / width - 0.5) / targetCam.zoom + targetCam.x;
        const wy = (sy / height - 0.5) / targetCam.zoom + targetCam.y;
        targetCam.zoom *= factor;
        clampCamObj(targetCam);
        targetCam.x = wx - (sx / width - 0.5) / targetCam.zoom;
        targetCam.y = wy - (sy / height - 0.5) / targetCam.zoom;
        clampCamObj(targetCam);
        kickSmooth();
    };

    const pickAt = (sx: number, sy: number): Pt | null => {
        let best: Pt | null = null;
        let bestD = 18 * 18;
        for (const p of points.value) {
            const s = project(p);
            const d = (s.x - sx) ** 2 + (s.y - sy) ** 2;
            if (d < bestD) {
                bestD = d;
                best = p;
            }
        }
        return best;
    };

    const onPointerDown = (e: PointerEvent) => {
        if (props.ambient) return;
        const el = wrapRef.value!;
        el.setPointerCapture(e.pointerId);
        pointers.set(e.pointerId, { x: e.offsetX, y: e.offsetY });
        if (pointers.size === 2) {
            const [a, b] = [...pointers.values()];
            pinchDist = Math.hypot(a.x - b.x, a.y - b.y);
        }
        dragging = true;
        moved = false;
        lastPtr = { x: e.offsetX, y: e.offsetY };
        vel.x = vel.y = 0;
    };

    const onPointerMove = (e: PointerEvent) => {
        if (props.ambient) return;
        if (pointers.has(e.pointerId)) pointers.set(e.pointerId, { x: e.offsetX, y: e.offsetY });

        if (pointers.size === 2) {
            const [a, b] = [...pointers.values()];
            const d = Math.hypot(a.x - b.x, a.y - b.y);
            if (pinchDist > 0 && d > 0) {
                zoomAt((a.x + b.x) / 2, (a.y + b.y) / 2, d / pinchDist);
                pinchDist = d;
            }
            return;
        }

        if (dragging && e.shiftKey && is3D.value) {
            const dy = e.offsetY - lastPtr.y;
            if (Math.abs(dy) > 1) moved = true;
            tilt.value = Math.min(1.25, Math.max(0, tilt.value - dy * 0.004));
            lastPtr = { x: e.offsetX, y: e.offsetY };
            scheduleRender();
        } else if (dragging) {
            const dx = e.offsetX - lastPtr.x;
            const dy = e.offsetY - lastPtr.y;
            if (Math.abs(dx) + Math.abs(dy) > 2) moved = true;
            const yDiv = is3D.value ? Math.max(0.35, Math.cos(tilt.value)) : 1;
            cam.x -= dx / width / cam.zoom;
            cam.y -= dy / height / cam.zoom / yDiv;
            targetCam.x = cam.x;
            targetCam.y = cam.y;
            clampCam();
            vel.x = dx / width;
            vel.y = dy / height;
            lastPtr = { x: e.offsetX, y: e.offsetY };
            scheduleRender();
        } else {
            const hit = pickAt(e.offsetX, e.offsetY);
            hoverPoint.value = hit;
            hoverScreen.value = { x: e.offsetX, y: e.offsetY };
            render();
        }
    };

    const onPointerUp = (e: PointerEvent) => {
        if (props.ambient) return;
        pointers.delete(e.pointerId);
        if (pointers.size < 2) pinchDist = 0;
        if (pointers.size === 0) dragging = false;
        if (!dragging && (Math.abs(vel.x) > 0.00002 || Math.abs(vel.y) > 0.00002)) kickInertia();
        if (!moved) {
            const hit = pickAt(e.offsetX, e.offsetY);
            if (hit) {
                const url = localePath(`/model/${hit.id}`);
                // Prefer a new tab; when the popup is blocked (embedded
                // browsers, strict blockers) fall back to same-tab
                // navigation so the click is never lost.
                const win = window.open(url, "_blank", "noopener");
                if (!win) navigateTo(url);
            }
        }
    };

    const onWheel = (e: WheelEvent) => {
        if (props.ambient) return;
        e.preventDefault();
        e.stopPropagation();
        const rect = wrapRef.value!.getBoundingClientRect();
        zoomAt(e.clientX - rect.left, e.clientY - rect.top, Math.pow(1.0022, -e.deltaY));
    };

    const resetView = () => {
        targetCam.x = 0.5;
        targetCam.y = 0.5;
        targetCam.zoom = 1;
        tilt.value = TILT_DEFAULT;
        vel.x = vel.y = 0;
        kickSmooth();
    };
    const set3D = (v: boolean) => {
        is3D.value = v;
        render();
    };
    const zoomIn = () => zoomAt(width / 2, height / 2, 1.6);
    const zoomOut = () => zoomAt(width / 2, height / 2, 1 / 1.6);

    defineExpose({ highlightProvider, resetView, zoomIn, zoomOut, set3D });

    // ---------- lifecycle ----------
    const resize = () => {
        const canvas = canvasRef.value;
        const el = wrapRef.value;
        if (!canvas || !el) return;
        // dpr 2 quadruples pixel count; lite mode trades retina sharpness
        // for fill-rate on low-spec machines.
        dpr = lite.value ? 1 : Math.min(2, window.devicePixelRatio || 1);
        const rect = el.getBoundingClientRect();
        width = rect.width;
        height = rect.height;
        canvas.width = Math.round(width * dpr);
        canvas.height = Math.round(height * dpr);
        const ctx = canvas.getContext("2d");
        ctx?.setTransform(dpr, 0, 0, dpr, 0, 0);
        if (!props.ambient) render();
    };

    let ro: ResizeObserver | null = null;
    let io: IntersectionObserver | null = null;

    // Quality toggle affects dpr (canvas re-init) and sprite baking —
    // re-run both, then the ambient loop / interactive render repaints.
    watch(lite, () => {
        sprites.clear();
        resize();
    });

    onMounted(() => {
        resize();
        ro = new ResizeObserver(() => {
            resize();
            sprites.clear();
        });
        ro.observe(wrapRef.value!);
        if (props.ambient) {
            start();
            io = new IntersectionObserver(([entry]) => (entry.isIntersecting ? start() : stop()), { threshold: 0.01 });
            io.observe(wrapRef.value!);
        } else {
            try {
                render();
            } catch (err) {
                console.error("[atlas] interactive mount render failed:", err);
            }
            // Survival repaints: early-page churn (font swap, layout
            // settle, theme init) can blank the canvas outside our
            // instrumented paths; repaint once the page is fully settled.
            const repaint = () => {
                try {
                    render();
                } catch {}
            };
            (document as any).fonts?.ready?.then(repaint);
            setTimeout(repaint, 2500);
            setTimeout(repaint, 4000);
        }

        // Pointer/wheel handlers are bound via template events (@pointerdown
        // etc.) — no manual addEventListener to lose.
    });

    onUnmounted(() => {
        cancelAnimationFrame(inertiaRaf);
        cancelAnimationFrame(smoothRaf);
        stop();
        ro?.disconnect();
        io?.disconnect();
    });

    function hash01(n: number): number {
        let x = Math.sin(n * 127.1) * 43758.5453;
        return x - Math.floor(x);
    }
</script>
