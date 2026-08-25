# chvorinov-t

A command-line tool that computes casting solidification times with
**Chvorinov's rule**. Given a casting volume `V`, surface area `A`, mold
constant `C`, and exponent `n`, it computes the solidification modulus
`M = V/A` and the freezing time `tf = C·Mⁿ`, then compares the moduli of
different shape classes (cube, plate, cylinder, sphere) at the same
volume.

This is a solidification *kernel*, not a foundry scheduling or MES system:
it answers the metallurgical question "how long does this casting take to
freeze" for a single casting at a time.

## What it computes

The freezing time of a casting is governed by its volume-to-surface ratio:

```
M  = V / A
tf = C · Mⁿ
```

For the standard shape classes the textbook moduli are exact or asymptotic
references:

```
cube     M = a/6     (edge a)
plate    M = t/2     (thickness t, infinite plate)
cylinder M = r/2     (radius r, infinite cylinder)
sphere   M = r/3     (radius r)
```

The tool also verifies the cross rules that follow from the model:

- doubling every linear dimension (similarity scaling) doubles `M`, so at
  `n = 2` the freezing time is quadrupled;
- doubling the mold constant `C` doubles `tf`;
- at equal volume a flatter plate exposes more surface, so it has a
  smaller `M` and a shorter `tf`;
- shapes that share the same `V/A` have the same `tf`, regardless of the
  shape label.

An optional riser check compares the riser modulus against the casting
modulus and warns when the riser would not solidify last. An optional
superheat term extends `tf` by the superheat-to-latent-heat ratio; with
zero superheat the result is the pure Chvorinov time.

## Units

The tool works in centimetre units: `V` in cm³, `A` in cm², `M` in cm,
`C` in min/cmⁿ, and `tf` in minutes. This matches the way Chvorinov's rule
is normally quoted in casting handbooks.

## Usage

The input is a JSON file describing the casting:

```json
{
  "label": "steel-cube",
  "shape": "cube",
  "edge": 10,
  "C": 3,
  "n": 2,
  "superheat_k": 0,
  "riser": {
    "shape": "cylinder",
    "radius": 6,
    "height": 30
  }
}
```

`shape` selects the geometry and its dimension fields (`cube`: `edge`;
`plate`: `thickness`/`width`/`length`; `cylinder`: `radius`/`height`;
`sphere`: `radius`). A `shape` of `custom` (or omitted) reads `volume` and
`area` directly. `C` is required; `n` defaults to 2.

```bash
# compute M and tf, with shape comparison and cross rules
go run . freeze example/steel-cube.json

# shape comparison at the volume of the input file
go run . compare example/steel-cube.json

# similarity-scaling and mold-constant cross rules
go run . scale example/steel-cube.json

# riser check
go run . riser example/steel-cube.json

# validate a file without computing
go run . validate example/steel-cube.json
```

The freezing time of the example cube comes out at 8.333 min:
`M = 10/6 = 1.667 cm`, `tf = 3 · (10/6)² = 8.333 min`.

### HTTP service

`chvorinov-t -http :8080` starts the calculation API; the request body is
the same casting JSON shown above:

- `GET /health` liveness probe;
- `POST /api/freeze` returns modulus, freezing time, shape table and riser check;
- `POST /api/compare` returns the equal-volume shape comparison;
- `POST /api/scale` returns the similarity-scaling cross rule;
- `POST /api/riser` returns the riser modulus check;
- `POST /api/validate` validates the casting input.

## Input validation

`V ≤ 0`, `A ≤ 0`, `C ≤ 0`, and `n ≤ 0` are rejected with an error naming
the offending quantity; invalid input exits non-zero with the message on
stderr. A surface area smaller than the enclosing sphere of the same
volume is geometrically impossible and is rejected as well.

## Build & test

```bash
go build ./...
go test ./...
```

## Layout

- `internal/geometry` — shape volumes, areas, moduli, and the physical
  feasibility bound (minimum sphere area).
- `internal/chvorinov` — the freezing-time kernel: validation, `tf`, riser
  check, superheat term, scaling and comparison cross rules.
- `internal/cli` — subcommand routing, JSON input, and report formatting.
- `example/` — ready-to-run casting JSON files.
