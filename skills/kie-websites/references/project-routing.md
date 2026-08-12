# Project routing

| Lane | Local responsibility | Kie responsibility |
| --- | --- | --- |
| Website | pages, components, SSR/SSG, forms, SEO, analytics, deploy | brand imagery, product media, hero video, illustrations |
| Media app | UI, server proxy, auth, persistence, quotas, moderation, job history | provider generation tasks and result URLs |
| Browser game | loop, input, physics, state, persistence, multiplayer, packaging | concept art, sprites/textures, backgrounds, sound/music references |

Security rules:

- keep Kie credentials server-side or in the local CLI credential store;
- validate uploads, MIME, size, and ownership;
- rate-limit live generation and make cost-bearing actions explicit;
- sanitize generated filenames and never execute model-supplied text;
- provide deletion/retention behavior for uploaded human likeness media;
- label generated media and avoid deceptive identity use.

Asset manifest fields: semantic ID, use, prompt/brief ID, model, generation/task ID, source URL, local optimized file, dimensions/duration, rights/provenance, approval state, and fallback.
