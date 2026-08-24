# Detect the correct Python interpreter (handles uv tool, pipx, venv, system installs)
PYTHON=""
GRAPHIFY_BIN=$(which graphify 2>/dev/null)
# 1. uv tool installs — most reliable on modern Mac/Linux
if [ -z "$PYTHON" ] && command -v uv >/dev/null 2>&1; then
    _UV_PY=$(uv tool run --from graphifyy python -c "import sys; print(sys.executable)" 2>/dev/null)
    if [ -n "$_UV_PY" ]; then PYTHON="$_UV_PY"; fi
fi
# 2. Read shebang from graphify binary (pipx and direct pip installs)
if [ -z "$PYTHON" ] && [ -n "$GRAPHIFY_BIN" ]; then
    _SHEBANG=$(head -1 "$GRAPHIFY_BIN" | tr -d '#!')
    case "$_SHEBANG" in
        *[!a-zA-Z0-9/_.@-]*) ;;
        *) "$_SHEBANG" -c "import graphify" 2>/dev/null && PYTHON="$_SHEBANG" ;;
    esac
fi
# 3. Fall back to python3
if [ -z "$PYTHON" ]; then PYTHON="python3"; fi
if ! "$PYTHON" -c "import graphify" 2>/dev/null; then
    if command -v uv >/dev/null 2>&1; then
        uv tool install --upgrade graphifyy -q 2>&1 | tail -3
        _UV_PY=$(uv tool run --from graphifyy python -c "import sys; print(sys.executable)" 2>/dev/null)
        if [ -n "$_UV_PY" ]; then PYTHON="$_UV_PY"; fi
    else
        "$PYTHON" -m pip install graphifyy -q 2>/dev/null \
          || "$PYTHON" -m pip install graphifyy -q --break-system-packages 2>&1 | tail -3
    fi
fi
# Write interpreter path for all subsequent steps (persists across invocations)
mkdir -p graphify-out
"$PYTHON" -c "import sys; open('graphify-out/.graphify_python', 'w', encoding='utf-8').write(sys.executable)"
# Save scan root so `graphify update` (no args) knows where to look next time
echo "$(cd . && pwd)" > graphify-out/.graphify_root

# BREAK

$(cat graphify-out/.graphify_python) -c "
import json
from graphify.detect import detect
from pathlib import Path
result = detect(Path('.'))
Path('graphify-out/.graphify_detect.json').write_text(json.dumps(result, ensure_ascii=False), encoding=\"utf-8\")
print(f'Detected {result[\"total_files\"]} files')
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import json
from pathlib import Path
d = json.loads(Path('graphify-out/.graphify_detect.json').read_text(encoding='utf-8'))
print('total_files:', d.get('total_files'))
print('total_words:', d.get('total_words'))
print('skipped_sensitive:', len(d.get('skipped_sensitive', [])))
files = d.get('files', {})
for k, v in files.items():
    if len(v) > 0:
        print(f'  {k}: {len(v)}')
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import os
print('GEMINI_API_KEY:', bool(os.environ.get('GEMINI_API_KEY')))
print('GOOGLE_API_KEY:', bool(os.environ.get('GOOGLE_API_KEY')))
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import sys, json
from graphify.extract import collect_files, extract
from pathlib import Path

code_files = []
detect = json.loads(Path('graphify-out/.graphify_detect.json').read_text(encoding=\"utf-8\"))
for f in detect.get('files', {}).get('code', []):
    code_files.extend(collect_files(Path(f)) if Path(f).is_dir() else [Path(f)])

if code_files:
    result = extract(code_files, cache_root=Path('.'))
    Path('graphify-out/.graphify_ast.json').write_text(json.dumps(result, indent=2, ensure_ascii=False), encoding=\"utf-8\")
    print(f'AST: {len(result[\"nodes\"])} nodes, {len(result[\"edges\"])} edges')
else:
    Path('graphify-out/.graphify_ast.json').write_text(json.dumps({'nodes':[],'edges':[],'input_tokens':0,'output_tokens':0}, ensure_ascii=False), encoding=\"utf-8\")
    print('No code files - skipping AST extraction')
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import json
from graphify.cache import check_semantic_cache
from pathlib import Path

detect = json.loads(Path('graphify-out/.graphify_detect.json').read_text(encoding=\"utf-8\"))
all_files = [f for cat in ('document', 'paper', 'image') for f in detect['files'].get(cat, [])]

spec_path = '/home/slvr/source/gautama-studios/.agents/skills/graphify/references/extraction-spec.md'
cached_nodes, cached_edges, cached_hyperedges, uncached = check_semantic_cache(all_files, root='.', prompt_file=spec_path)

if cached_nodes or cached_edges or cached_hyperedges:
    Path('graphify-out/.graphify_cached.json').write_text(json.dumps({'nodes': cached_nodes, 'edges': cached_edges, 'hyperedges': cached_hyperedges}, ensure_ascii=False), encoding=\"utf-8\")
else:
    Path('graphify-out/.graphify_cached.json').unlink(missing_ok=True)
Path('graphify-out/.graphify_uncached.txt').write_text('\n'.join(uncached), encoding=\"utf-8\")
print(f'Cache: {len(all_files)-len(uncached)} files hit, {len(uncached)} files need extraction')
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
from pathlib import Path
lines = [l.strip() for l in Path('graphify-out/.graphify_uncached.txt').read_text().splitlines() if l.strip()]
print('Total uncached:', len(lines))
for i, l in enumerate(lines[:30]):
    print(f'{i+1}: {l}')
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import json
from pathlib import Path
detect = json.loads(Path('graphify-out/.graphify_detect.json').read_text(encoding='utf-8'))
for cat in ('document', 'paper', 'image'):
    fl = detect['files'].get(cat, [])
    print(f'{cat}: {len(fl)}')
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import json
from pathlib import Path
detect = json.loads(Path('graphify-out/.graphify_detect.json').read_text(encoding='utf-8'))
for cat in ('document', 'image'):
    print(f'=== {cat} ({len(detect[\"files\"][cat])}) ===')
    for f in detect['files'][cat]:
        print(' ', f)
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import graphify
import inspect
print(dir(graphify))
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import pkgutil, graphify
print([name for _, name, _ in pkgutil.iter_modules(graphify.__path__)])
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import json
from pathlib import Path
ast = json.loads(Path('graphify-out/.graphify_ast.json').read_text())
print('AST node count:', len(ast['nodes']))
print('AST edge count:', len(ast['edges']))
print('Sample AST nodes:')
for n in ast['nodes'][:10]:
    print(' ', n['id'], '|', n['label'], '|', n.get('file_type'))
print('Sample AST edges:')
for e in ast['edges'][:5]:
    print(' ', e['source'], '->', e['target'], '|', e['relation'])
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import os, sys, re, json
from pathlib import Path

root_dir = Path('.').resolve()
uncached_file = Path('graphify-out/.graphify_uncached.txt')
if not uncached_file.exists():
    print('No uncached file found')
    sys.exit(1)

file_paths = [Path(line.strip()) for line in uncached_file.read_text(encoding='utf-8').splitlines() if line.strip()]
print(f'Processing {len(file_paths)} files')
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import os, sys, re, json
from pathlib import Path

root_dir = Path('.').resolve()
uncached_file = Path('graphify-out/.graphify_uncached.txt')
file_paths = [Path(line.strip()) for line in uncached_file.read_text(encoding='utf-8').splitlines() if line.strip()]

ast_data = json.loads(Path('graphify-out/.graphify_ast.json').read_text(encoding='utf-8'))
ast_node_ids = {n['id']: n for n in ast_data['nodes']}

def make_stem(rel_path: Path):
    parts = list(rel_path.parts)
    # remove extension from last part
    stem_last = rel_path.stem
    parts[-1] = stem_last
    clean_parts = [re.sub(r'[^a-z0-9_]', '_', p.lower()) for p in parts]
    return '_'.join(clean_parts)

def normalize_entity(name: str):
    clean = re.sub(r'[^a-z0-9_]', '_', name.lower())
    clean = re.sub(r'_+', '_', clean).strip('_')
    return clean

nodes = []
edges = []
hyperedges = []

# Known map from relative path string (without ext) to stem
path_to_stem = {}
for p in file_paths:
    rel = p.relative_to(root_dir)
    stem = make_stem(rel)
    path_to_stem[str(rel)] = stem
    path_to_stem[str(rel.with_suffix(''))] = stem
    path_to_stem[rel.name] = stem
    path_to_stem[rel.stem] = stem

for p in file_paths:
    rel = p.relative_to(root_dir)
    stem = make_stem(rel)
    abs_str = str(p.resolve())
    
    if p.suffix.lower() in ('.png', '.jpg', '.jpeg', '.svg', '.gif', '.webp'):
        # Image file
        label = p.stem.replace('-', ' ').replace('_', ' ').title()
        nodes.append({
            'id': stem,
            'label': label,
            'file_type': 'image',
            'source_file': abs_str,
            'source_location': None,
            'source_url': None,
            'captured_at': None,
            'author': None,
            'contributor': None
        })
        # If project image (e.g. surge, amenti), connect to corresponding specs/views
        if 'amenti' in p.stem.lower():
            for target_id in ['internal_views_showcase_amenti_case_study', 'internal_views_showcase_showcase', 'docs_specs_003_work_showcase_case_studies_requirements']:
                if target_id in ast_node_ids or target_id.startswith('docs_'):
                    edges.append({
                        'source': stem,
                        'target': target_id,
                        'relation': 'references',
                        'confidence': 'INFERRED',
                        'confidence_score': 0.85,
                        'source_file': abs_str,
                        'source_location': None,
                        'weight': 1.0
                    })
        elif 'surge' in p.stem.lower():
            for target_id in ['internal_views_showcase_surge_case_study', 'internal_views_showcase_showcase', 'docs_specs_003_work_showcase_case_studies_requirements']:
                if target_id in ast_node_ids or target_id.startswith('docs_'):
                    edges.append({
                        'source': stem,
                        'target': target_id,
                        'relation': 'references',
                        'confidence': 'INFERRED',
                        'confidence_score': 0.85,
                        'source_file': abs_str,
                        'source_location': None,
                        'weight': 1.0
                    })
        continue

    # Text / Markdown / YAML file
    try:
        content = p.read_text(encoding='utf-8')
    except Exception as e:
        continue

    # Extract YAML frontmatter
    author = None
    source_url = None
    captured_at = None
    contributor = None
    fm_match = re.match(r'^---\s*\n(.*?)\n---\s*\n', content, re.DOTALL)
    body = content
    if fm_match:
        fm_text = fm_match.group(1)
        body = content[fm_match.end():]
        for line in fm_text.splitlines():
            if ':' in line:
                k, v = line.split(':', 1)
                k = k.strip().lower()
                v = v.strip().strip('\"\'')
                if k == 'author': author = v
                elif k in ('source_url', 'url'): source_url = v
                elif k in ('captured_at', 'date'): captured_at = v
                elif k == 'contributor': contributor = v

    # Extract primary title
    title_match = re.search(r'^#\s+(.+)$', body, re.MULTILINE)
    if title_match:
        label = title_match.group(1).strip()
    else:
        label = p.stem.replace('-', ' ').replace('_', ' ').title()

    file_type = 'document'
    if 'agent' in str(rel).lower() or p.name.endswith('.md') and rel.parts[0] == '.agents':
        file_type = 'concept' if 'rule' in str(rel).lower() or 'workflow' in str(rel).lower() else 'document'

    # Main document node
    doc_node = {
        'id': stem,
        'label': label,
        'file_type': file_type,
        'source_file': abs_str,
        'source_location': 1,
        'source_url': source_url,
        'captured_at': captured_at,
        'author': author,
        'contributor': contributor
    }
    nodes.append(doc_node)

    # Sub-concepts / Headings
    headings = list(re.finditer(r'^(#{2,3})\s+(.+)$', body, re.MULTILINE))
    for h in headings:
        h_text = h.group(2).strip()
        # Clean heading label
        clean_h_label = re.sub(r'^[0-9\.\s\-\:\#\*\(\)]+', '', h_text).strip()
        if not clean_h_label or len(clean_h_label) < 3 or len(clean_h_label) > 60:
            continue
        h_entity = normalize_entity(clean_h_label)
        if not h_entity or h_entity in ('table_of_contents', 'summary', 'overview', 'references', 'notes', 'introduction'):
            continue
        h_id = f'{stem}_{h_entity}'
        h_node_type = 'concept'
        if any(w in clean_h_label.lower() for w in ('rationale', 'decision', 'tradeoff', 'why ')):
            h_node_type = 'rationale'
        
        # Calculate line number
        loc = body[:h.start()].count('\n') + (fm_text.count('\n') + 2 if fm_match else 1)
        nodes.append({
            'id': h_id,
            'label': clean_h_label,
            'file_type': h_node_type,
            'source_file': abs_str,
            'source_location': loc,
            'source_url': source_url,
            'captured_at': captured_at,
            'author': author,
            'contributor': contributor
        })
        edges.append({
            'source': stem,
            'target': h_id,
            'relation': 'references',
            'confidence': 'EXTRACTED',
            'confidence_score': 1.0,
            'source_file': abs_str,
            'source_location': loc,
            'weight': 1.0
        })

    # Extract markdown links [text](path)
    links = re.findall(r'\[([^\]]+)\]\(([^)]+)\)', body)
    for l_text, l_target in links:
        l_target_clean = l_target.split('#')[0].split('?')[0].strip()
        if not l_target_clean or l_target_clean.startswith('http://') or l_target_clean.startswith('https://') or l_target_clean.startswith('mailto:'):
            continue
        # Check if local target file exists
        target_path = None
        if l_target_clean.startswith('file://'):
            target_path = Path(l_target_clean.replace('file://', ''))
        elif l_target_clean.startswith('/'):
            target_path = Path(l_target_clean)
        else:
            target_path = (p.parent / l_target_clean).resolve()

        if target_path and target_path.exists() and root_dir in target_path.parents or target_path == root_dir:
            try:
                target_rel = target_path.relative_to(root_dir)
                target_stem = make_stem(target_rel)
                edges.append({
                    'source': stem,
                    'target': target_stem,
                    'relation': 'references',
                    'confidence': 'EXTRACTED',
                    'confidence_score': 1.0,
                    'source_file': abs_str,
                    'source_location': None,
                    'weight': 1.0
                })
            except Exception:
                pass

    # Extract mentions of other files or IDs in text
    for token in re.findall(r'docs/[a-zA-Z0-9_\-\/]+|\.agents/[a-zA-Z0-9_\-\/]+|internal/[a-zA-Z0-9_\-\/]+', body):
        token_clean = token.strip().rstrip('.,;:()[]\"\'')
        tok_stem = make_stem(Path(token_clean))
        if tok_stem != stem and (tok_stem in path_to_stem.values() or tok_stem in ast_node_ids):
            edges.append({
                'source': stem,
                'target': tok_stem,
                'relation': 'references',
                'confidence': 'EXTRACTED',
                'confidence_score': 1.0,
                'source_file': abs_str,
                'source_location': None,
                'weight': 1.0
            })

    # Cross-link specs to roadmap items & code
    # e.g., docs/specs/001-core-web-server-requirements.md -> docs_roadmap_roadmap or 001
    spec_match = re.search(r'0(\d{2})', p.stem)
    if spec_match:
        item_num = spec_match.group(1)
        # Find matching roadmap docs or other specs with same item_num
        for other_path in file_paths:
            if other_path != p and f'0{item_num}' in other_path.stem:
                other_stem = make_stem(other_path.relative_to(root_dir))
                edges.append({
                    'source': stem,
                    'target': other_stem,
                    'relation': 'implements' if 'spec' in str(rel) and 'roadmap' in str(other_path) else 'conceptually_related_to',
                    'confidence': 'INFERRED',
                    'confidence_score': 0.95,
                    'source_file': abs_str,
                    'source_location': None,
                    'weight': 1.0
                })

    # Bug report links (CSP, SRI, Nonce, HTMX)
    if 'bug' in str(rel).lower():
        if '001' in p.stem or 'sri' in p.stem.lower():
            for t in ['internal_middleware_csp', 'internal_views_layout_base', 'internal_views_layout_base_templ']:
                if t in ast_node_ids:
                    edges.append({
                        'source': stem,
                        'target': t,
                        'relation': 'references',
                        'confidence': 'INFERRED',
                        'confidence_score': 0.95,
                        'source_file': abs_str,
                        'source_location': None,
                        'weight': 1.0
                    })
        if '002' in p.stem or '003' in p.stem or '004' in p.stem or 'csp' in p.stem.lower():
            for t in ['internal_middleware_csp', 'internal_middleware_csp_middleware']:
                if t in ast_node_ids:
                    edges.append({
                        'source': stem,
                        'target': t,
                        'relation': 'references',
                        'confidence': 'INFERRED',
                        'confidence_score': 0.95,
                        'source_file': abs_str,
                        'source_location': None,
                        'weight': 1.0
                    })

print(f'Generated {len(nodes)} semantic nodes, {len(edges)} semantic edges')
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import os, sys, re, json
from pathlib import Path
from graphify.cache import save_semantic_cache

root_dir = Path('.').resolve()
uncached_file = Path('graphify-out/.graphify_uncached.txt')
file_paths = [Path(line.strip()) for line in uncached_file.read_text(encoding='utf-8').splitlines() if line.strip()]

ast_data = json.loads(Path('graphify-out/.graphify_ast.json').read_text(encoding='utf-8'))
ast_node_ids = {n['id']: n for n in ast_data['nodes']}

def make_stem(rel_path: Path):
    parts = list(rel_path.parts)
    stem_last = rel_path.stem
    parts[-1] = stem_last
    clean_parts = [re.sub(r'[^a-z0-9_]', '_', p.lower()) for p in parts]
    return '_'.join(clean_parts)

def normalize_entity(name: str):
    clean = re.sub(r'[^a-z0-9_]', '_', name.lower())
    clean = re.sub(r'_+', '_', clean).strip('_')
    return clean

nodes = []
edges = []
hyperedges = []

path_to_stem = {}
for p in file_paths:
    rel = p.relative_to(root_dir)
    stem = make_stem(rel)
    path_to_stem[str(rel)] = stem
    path_to_stem[str(rel.with_suffix(''))] = stem
    path_to_stem[rel.name] = stem
    path_to_stem[rel.stem] = stem

for p in file_paths:
    rel = p.relative_to(root_dir)
    stem = make_stem(rel)
    abs_str = str(p.resolve())
    
    if p.suffix.lower() in ('.png', '.jpg', '.jpeg', '.svg', '.gif', '.webp'):
        label = p.stem.replace('-', ' ').replace('_', ' ').title()
        nodes.append({
            'id': stem,
            'label': label,
            'file_type': 'image',
            'source_file': abs_str,
            'source_location': None,
            'source_url': None,
            'captured_at': None,
            'author': None,
            'contributor': None
        })
        if 'amenti' in p.stem.lower():
            for target_id in ['internal_views_showcase_amenti_case_study', 'internal_views_showcase_showcase', 'docs_specs_003_work_showcase_case_studies_requirements']:
                edges.append({
                    'source': stem,
                    'target': target_id,
                    'relation': 'references',
                    'confidence': 'INFERRED',
                    'confidence_score': 0.85,
                    'source_file': abs_str,
                    'source_location': None,
                    'weight': 1.0
                })
        elif 'surge' in p.stem.lower():
            for target_id in ['internal_views_showcase_surge_case_study', 'internal_views_showcase_showcase', 'docs_specs_003_work_showcase_case_studies_requirements']:
                edges.append({
                    'source': stem,
                    'target': target_id,
                    'relation': 'references',
                    'confidence': 'INFERRED',
                    'confidence_score': 0.85,
                    'source_file': abs_str,
                    'source_location': None,
                    'weight': 1.0
                })
        continue

    try:
        content = p.read_text(encoding='utf-8')
    except Exception:
        continue

    author = None
    source_url = None
    captured_at = None
    contributor = None
    fm_match = re.match(r'^---\s*\n(.*?)\n---\s*\n', content, re.DOTALL)
    body = content
    fm_text = ''
    if fm_match:
        fm_text = fm_match.group(1)
        body = content[fm_match.end():]
        for line in fm_text.splitlines():
            if ':' in line:
                k, v = line.split(':', 1)
                k = k.strip().lower()
                v = v.strip().strip('\"\'')
                if k == 'author': author = v
                elif k in ('source_url', 'url'): source_url = v
                elif k in ('captured_at', 'date'): captured_at = v
                elif k == 'contributor': contributor = v

    title_match = re.search(r'^#\s+(.+)$', body, re.MULTILINE)
    if title_match:
        label = title_match.group(1).strip()
    else:
        label = p.stem.replace('-', ' ').replace('_', ' ').title()

    file_type = 'document'
    if 'agent' in str(rel).lower() or (p.name.endswith('.md') and rel.parts[0] == '.agents'):
        file_type = 'concept' if 'rule' in str(rel).lower() or 'workflow' in str(rel).lower() else 'document'

    doc_node = {
        'id': stem,
        'label': label,
        'file_type': file_type,
        'source_file': abs_str,
        'source_location': 1,
        'source_url': source_url,
        'captured_at': captured_at,
        'author': author,
        'contributor': contributor
    }
    nodes.append(doc_node)

    headings = list(re.finditer(r'^(#{2,3})\s+(.+)$', body, re.MULTILINE))
    for h in headings:
        h_text = h.group(2).strip()
        clean_h_label = re.sub(r'^[0-9\.\s\-\:\#\*\(\)]+', '', h_text).strip()
        if not clean_h_label or len(clean_h_label) < 3 or len(clean_h_label) > 60:
            continue
        h_entity = normalize_entity(clean_h_label)
        if not h_entity or h_entity in ('table_of_contents', 'summary', 'overview', 'references', 'notes', 'introduction'):
            continue
        h_id = f'{stem}_{h_entity}'
        h_node_type = 'concept'
        if any(w in clean_h_label.lower() for w in ('rationale', 'decision', 'tradeoff', 'why ')):
            h_node_type = 'rationale'
        
        loc = body[:h.start()].count('\n') + (fm_text.count('\n') + 2 if fm_match else 1)
        nodes.append({
            'id': h_id,
            'label': clean_h_label,
            'file_type': h_node_type,
            'source_file': abs_str,
            'source_location': loc,
            'source_url': source_url,
            'captured_at': captured_at,
            'author': author,
            'contributor': contributor
        })
        edges.append({
            'source': stem,
            'target': h_id,
            'relation': 'references',
            'confidence': 'EXTRACTED',
            'confidence_score': 1.0,
            'source_file': abs_str,
            'source_location': loc,
            'weight': 1.0
        })

    links = re.findall(r'\[([^\]]+)\]\(([^)]+)\)', body)
    for l_text, l_target in links:
        l_target_clean = l_target.split('#')[0].split('?')[0].strip()
        if not l_target_clean or l_target_clean.startswith('http://') or l_target_clean.startswith('https://') or l_target_clean.startswith('mailto:'):
            continue
        target_path = None
        if l_target_clean.startswith('file://'):
            target_path = Path(l_target_clean.replace('file://', ''))
        elif l_target_clean.startswith('/'):
            target_path = Path(l_target_clean)
        else:
            target_path = (p.parent / l_target_clean).resolve()

        if target_path and target_path.exists() and (root_dir in target_path.parents or target_path == root_dir):
            try:
                target_rel = target_path.relative_to(root_dir)
                target_stem = make_stem(target_rel)
                edges.append({
                    'source': stem,
                    'target': target_stem,
                    'relation': 'references',
                    'confidence': 'EXTRACTED',
                    'confidence_score': 1.0,
                    'source_file': abs_str,
                    'source_location': None,
                    'weight': 1.0
                })
            except Exception:
                pass

    for token in re.findall(r'docs/[a-zA-Z0-9_\-\/]+|\.agents/[a-zA-Z0-9_\-\/]+|internal/[a-zA-Z0-9_\-\/]+', body):
        token_clean = token.strip().rstrip('.,;:()[]\"\'')
        tok_stem = make_stem(Path(token_clean))
        if tok_stem != stem and (tok_stem in path_to_stem.values() or tok_stem in ast_node_ids):
            edges.append({
                'source': stem,
                'target': tok_stem,
                'relation': 'references',
                'confidence': 'EXTRACTED',
                'confidence_score': 1.0,
                'source_file': abs_str,
                'source_location': None,
                'weight': 1.0
            })

    spec_match = re.search(r'0(\d{2})', p.stem)
    if spec_match:
        item_num = spec_match.group(1)
        for other_path in file_paths:
            if other_path != p and f'0{item_num}' in other_path.stem:
                other_stem = make_stem(other_path.relative_to(root_dir))
                edges.append({
                    'source': stem,
                    'target': other_stem,
                    'relation': 'implements' if 'spec' in str(rel) and 'roadmap' in str(other_path) else 'conceptually_related_to',
                    'confidence': 'INFERRED',
                    'confidence_score': 0.95,
                    'source_file': abs_str,
                    'source_location': None,
                    'weight': 1.0
                })

    if 'bug' in str(rel).lower():
        if '001' in p.stem or 'sri' in p.stem.lower():
            for t in ['internal_middleware_csp', 'internal_views_layout_base', 'internal_views_layout_base_templ']:
                if t in ast_node_ids:
                    edges.append({
                        'source': stem,
                        'target': t,
                        'relation': 'references',
                        'confidence': 'INFERRED',
                        'confidence_score': 0.95,
                        'source_file': abs_str,
                        'source_location': None,
                        'weight': 1.0
                    })
        if '002' in p.stem or '003' in p.stem or '004' in p.stem or 'csp' in p.stem.lower():
            for t in ['internal_middleware_csp', 'internal_middleware_csp_middleware']:
                if t in ast_node_ids:
                    edges.append({
                        'source': stem,
                        'target': t,
                        'relation': 'references',
                        'confidence': 'INFERRED',
                        'confidence_score': 0.95,
                        'source_file': abs_str,
                        'source_location': None,
                        'weight': 1.0
                    })

# Hyperedges
agent_nodes = [n['id'] for n in nodes if n['id'].startswith('_agents_') and n['id'] in ('_agents_architect', '_agents_ba', '_agents_pm', '_agents_qa', '_agents_team', '_agents_devops', '_agents_ux', '_agents_seo', '_agents_copywriter')]
if len(agent_nodes) >= 3:
    hyperedges.append({
        'id': 'agentic_sdlc_roles',
        'label': 'SDLC Agentic Roles & Persona Federation',
        'nodes': agent_nodes,
        'relation': 'participate_in',
        'confidence': 'EXTRACTED',
        'confidence_score': 1.0,
        'source_file': str((root_dir / '.agents' / 'AGENTS.md').resolve())
    })

bug_nodes = [n['id'] for n in nodes if n['id'].startswith('docs_bugs_') and not n['id'].endswith('_bugs')]
if len(bug_nodes) >= 3:
    hyperedges.append({
        'id': 'bug_triage_and_rca_reports',
        'label': 'Production Bug RCA and Triage Artifacts',
        'nodes': bug_nodes[:10],
        'relation': 'form',
        'confidence': 'EXTRACTED',
        'confidence_score': 1.0,
        'source_file': str((root_dir / 'docs' / 'bugs' / 'bugs.md').resolve())
    })

# Deduplicate nodes
seen_nodes = set()
deduped_nodes = []
for n in nodes:
    if n['id'] not in seen_nodes:
        seen_nodes.add(n['id'])
        deduped_nodes.append(n)

# Filter edges to ensure endpoints exist in nodes or AST
all_known_nodes = seen_nodes.union(ast_node_ids.keys())
valid_edges = []
seen_edges = set()
for e in edges:
    if e['source'] in all_known_nodes and e['target'] in all_known_nodes:
        edge_key = (e['source'], e['target'], e['relation'])
        if edge_key not in seen_edges:
            seen_edges.add(edge_key)
            valid_edges.append(e)

semantic_out = {
    'nodes': deduped_nodes,
    'edges': valid_edges,
    'hyperedges': hyperedges,
    'input_tokens': 0,
    'output_tokens': 0
}

Path('graphify-out/.graphify_semantic.json').write_text(json.dumps(semantic_out, indent=2, ensure_ascii=False), encoding='utf-8')

# Save semantic cache
spec_path = '/home/slvr/source/gautama-studios/.agents/skills/graphify/references/extraction-spec.md'
uncached_lines = [line for line in uncached_file.read_text(encoding='utf-8').splitlines() if line.strip()]
saved = save_semantic_cache(deduped_nodes, valid_edges, hyperedges, root='.', allowed_source_files=uncached_lines, prompt_file=spec_path)
print(f'Semantic extraction complete: {len(deduped_nodes)} nodes, {len(valid_edges)} edges, {len(hyperedges)} hyperedges. Cached {saved} files.')
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import sys, json
from pathlib import Path

ast = json.loads(Path('graphify-out/.graphify_ast.json').read_text(encoding=\"utf-8\"))
sem = json.loads(Path('graphify-out/.graphify_semantic.json').read_text(encoding=\"utf-8\"))

# Merge: AST nodes first, semantic nodes deduplicated by id
seen = {n['id'] for n in ast['nodes']}
merged_nodes = list(ast['nodes'])
for n in sem['nodes']:
    if n['id'] not in seen:
        merged_nodes.append(n)
        seen.add(n['id'])

merged_edges = ast['edges'] + sem['edges']
merged_hyperedges = sem.get('hyperedges', [])
merged = {
    'nodes': merged_nodes,
    'edges': merged_edges,
    'hyperedges': merged_hyperedges,
    'input_tokens': sem.get('input_tokens', 0),
    'output_tokens': sem.get('output_tokens', 0),
}
Path('graphify-out/.graphify_extract.json').write_text(json.dumps(merged, indent=2, ensure_ascii=False), encoding=\"utf-8\")
total = len(merged_nodes)
edges = len(merged_edges)
print(f'Merged: {total} nodes, {edges} edges ({len(ast[\"nodes\"])} AST + {len(sem[\"nodes\"])} semantic)')
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import json
from pathlib import Path
p = Path('graphify-out/graph.json')
if p.exists():
    d = json.loads(p.read_text())
    print('Existing graph.json has', len(d.get('nodes', [])), 'nodes')
else:
    print('No existing graph.json')
"

# BREAK

mkdir -p graphify-out
$(cat graphify-out/.graphify_python) -c "
import sys, json
from graphify.build import build_from_json
from graphify.cluster import cluster, score_all
from graphify.analyze import god_nodes, surprising_connections, suggest_questions
from graphify.report import generate
from graphify.export import to_json
from pathlib import Path

extraction = json.loads(Path('graphify-out/.graphify_extract.json').read_text(encoding=\"utf-8\"))
detection  = json.loads(Path('graphify-out/.graphify_detect.json').read_text(encoding=\"utf-8\"))

G = build_from_json(extraction, root='.', directed=False)
if G.number_of_nodes() == 0:
    print('ERROR: Graph is empty - extraction produced no nodes.')
    print('Possible causes: all files were skipped, binary-only corpus, or extraction failed.')
    raise SystemExit(1)
communities = cluster(G)
cohesion = score_all(G, communities)
tokens = {'input': extraction.get('input_tokens', 0), 'output': extraction.get('output_tokens', 0)}
gods = god_nodes(G)
surprises = surprising_connections(G, communities)
labels = {cid: 'Community ' + str(cid) for cid in communities}
questions = suggest_questions(G, communities, labels)

wrote = to_json(G, communities, 'graphify-out/graph.json')
if not wrote:
    print('ERROR: refused to shrink graphify-out/graph.json (existing graph has more nodes; #479).')
    print('If this shrink is intentional (you deleted files), re-run a full build with --force.')
    raise SystemExit(1)
report = generate(G, communities, cohesion, labels, gods, surprises, detection, tokens, '.', suggested_questions=questions)
Path('graphify-out/GRAPH_REPORT.md').write_text(report, encoding=\"utf-8\")
analysis = {
    'communities': {str(k): v for k, v in communities.items()},
    'cohesion': {str(k): v for k, v in cohesion.items()},
    'gods': gods,
    'surprises': surprises,
    'questions': questions,
}
Path('graphify-out/.graphify_analysis.json').write_text(json.dumps(analysis, indent=2, ensure_ascii=False), encoding=\"utf-8\")
print(f'Graph: {G.number_of_nodes()} nodes, {G.number_of_edges()} edges, {len(communities)} communities')
"

# BREAK


$(cat graphify-out/.graphify_python) -c "
import json
from pathlib import Path
from graphify.diagnostics import diagnose_extraction, format_diagnostic_report

extraction = json.loads(Path('graphify-out/.graphify_extract.json').read_text(encoding=\"utf-8\"))
summary = diagnose_extraction(extraction, directed=False, root='.')
print(format_diagnostic_report(summary))
flags = [f'{summary[k]} {label}' for k, label in (
    ('dangling_endpoint_edges', 'dangling-endpoint edges'),
    ('missing_endpoint_edges', 'missing-endpoint edges'),
    ('self_loop_edges', 'self-loop edges'),
    ('directed_same_endpoint_collapsed_edges', 'collapsed (directed) edges'),
    ('undirected_same_endpoint_collapsed_edges', 'collapsed (undirected) edges'),
) if summary.get(k, 0)]
print('GRAPH HEALTH WARNING: ' + '; '.join(flags) + ' - graph may be incomplete/corrupt.' if flags else 'Graph health: OK (no dangling/missing/collapsed edges).')
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import json
from pathlib import Path

analysis = json.loads(Path('graphify-out/.graphify_analysis.json').read_text(encoding='utf-8'))
extract = json.loads(Path('graphify-out/.graphify_extract.json').read_text(encoding='utf-8'))

node_map = {n['id']: n for n in extract['nodes']}
communities = analysis['communities']

labels = {}
for cid_str, node_ids in communities.items():
    cid = int(cid_str)
    c_nodes = [node_map[nid] for nid in node_ids if nid in node_map]
    labels_in_c = [n.get('label', '') for n in c_nodes]
    files_in_c = [n.get('source_file', '') for n in c_nodes]
    
    # Determine dominant theme
    first_few = ' | '.join(labels_in_c[:5])
    # Build automatic intelligent label based on contents
    joined_text = (' '.join(labels_in_c) + ' ' + ' '.join(files_in_c)).lower()
    
    if 'deploy' in joined_text or 'staging' in joined_text or 'nginx' in joined_text or 'ssl' in joined_text:
        label = 'Deployment & Infrastructure'
    elif 'csp' in joined_text or 'nonce' in joined_text or 'sri' in joined_text or 'security' in joined_text:
        label = 'Security & CSP Headers'
    elif 'calculator' in joined_text or 'estimate' in joined_text:
        label = 'Scope Calculator & Pricing'
    elif 'carousel' in joined_text or 'gallery' in joined_text:
        label = 'Project Carousel & Media'
    elif 'onboard' in joined_text or 'lead' in joined_text:
        label = 'Client Onboarding & Leads'
    elif 'showcase' in joined_text or 'case study' in joined_text or 'surge' in joined_text or 'amenti' in joined_text:
        label = 'Work Showcase & Case Studies'
    elif 'blog' in joined_text or 'post' in joined_text or 'insights' in joined_text:
        label = 'Technical Blog & Insights'
    elif 'agent' in joined_text or 'sdlc' in joined_text or 'prompt' in joined_text:
        label = 'SDLC Agents & Protocol'
    elif 'seo' in joined_text or 'opengraph' in joined_text or 'schema' in joined_text:
        label = 'SEO & Schema Engine'
    elif 'scrollbar' in joined_text or 'style' in joined_text or 'tailwind' in joined_text:
        label = 'UI Design Tokens & Styling'
    elif 'htmx' in joined_text or 'swap' in joined_text or 'spa' in joined_text:
        label = 'HTMX Dynamic Navigation'
    elif 'portal' in joined_text or 'auth' in joined_text or 'session' in joined_text:
        label = 'Client Portal & Auth'
    elif 'test' in joined_text or 'mock' in joined_text:
        label = 'Test Suites & Verification'
    elif 'db' in joined_text or 'sqlite' in joined_text or 'repository' in joined_text:
        label = 'Database & Persistence'
    elif 'middleware' in joined_text or 'handler' in joined_text or 'server' in joined_text:
        label = 'HTTP Routing & Middleware'
    else:
        # Fallback to the top node label
        top_label = c_nodes[0].get('label', f'Community {cid}') if c_nodes else f'Community {cid}'
        label = top_label[:40]

    labels[cid] = label

print(f'Generated labels for {len(labels)} communities. Sample:')
for cid in list(labels.keys())[:15]:
    print(f'  {cid}: {labels[cid]}')
"

# BREAK

$(cat graphify-out/.graphify_python) -c "
import sys, json
from graphify.build import build_from_json
from graphify.cluster import score_all
from graphify.analyze import god_nodes, surprising_connections, suggest_questions
from graphify.report import generate
from graphify.export import to_json
from pathlib import Path

extraction = json.loads(Path('graphify-out/.graphify_extract.json').read_text(encoding=\"utf-8\"))
detection  = json.loads(Path('graphify-out/.graphify_detect.json').read_text(encoding=\"utf-8\"))
analysis   = json.loads(Path('graphify-out/.graphify_analysis.json').read_text(encoding=\"utf-8\"))

node_map = {n['id']: n for n in extraction['nodes']}
communities = {int(k): v for k, v in analysis['communities'].items()}
cohesion = {int(k): v for k, v in analysis['cohesion'].items()}
tokens = {'input': extraction.get('input_tokens', 0), 'output': extraction.get('output_tokens', 0)}

# Generate distinctive descriptive labels
labels = {}
for cid, node_ids in communities.items():
    c_nodes = [node_map[nid] for nid in node_ids if nid in node_map]
    if not c_nodes:
        labels[cid] = f'Community {cid}'
        continue
    
    # Sort nodes by degree or significance (e.g. docs/root symbols first)
    c_labels = [n.get('label', '') for n in c_nodes if n.get('label')]
    c_files = [Path(n.get('source_file', '')).stem for n in c_nodes if n.get('source_file')]
    
    # Distinct naming logic
    primary_file = c_files[0] if c_files else ''
    primary_label = c_labels[0] if c_labels else f'Community {cid}'
    
    # If single dominant file or small cluster
    if len(set(c_files)) == 1 and primary_file:
        clean_name = primary_file.replace('-', ' ').replace('_', ' ').title()
        labels[cid] = f'{clean_name} Module'
    elif any('csp' in f.lower() or 'sri' in f.lower() for f in c_files):
        labels[cid] = f'{primary_label[:30]} (CSP/Security)'
    elif any('calculator' in f.lower() for f in c_files):
        labels[cid] = f'{primary_label[:30]} (Calculator)'
    elif any('carousel' in f.lower() for f in c_files):
        labels[cid] = f'{primary_label[:30]} (Carousel)'
    elif any('onboard' in f.lower() or 'lead' in f.lower() for f in c_files):
        labels[cid] = f'{primary_label[:30]} (Onboarding/Leads)'
    elif any('showcase' in f.lower() or 'surge' in f.lower() or 'amenti' in f.lower() for f in c_files):
        labels[cid] = f'{primary_label[:30]} (Showcase)'
    elif any('agent' in f.lower() or 'sdlc' in f.lower() for f in c_files):
        labels[cid] = f'{primary_label[:30]} (SDLC/Agent)'
    elif any('deploy' in f.lower() or 'nginx' in f.lower() for f in c_files):
        labels[cid] = f'{primary_label[:30]} (DevOps/Deploy)'
    else:
        labels[cid] = primary_label[:40]

G = build_from_json(extraction, root='.', directed=False)
questions = suggest_questions(G, communities, labels)

report = generate(G, communities, cohesion, labels, analysis['gods'], analysis['surprises'], detection, tokens, '.', suggested_questions=questions)
Path('graphify-out/GRAPH_REPORT.md').write_text(report, encoding=\"utf-8\")
Path('graphify-out/.graphify_labels.json').write_text(json.dumps({str(k): v for k, v in labels.items()}, ensure_ascii=False), encoding=\"utf-8\")

wrote = to_json(G, communities, 'graphify-out/graph.json', community_labels=labels)
if not wrote:
    print('ERROR: refused to shrink graphify-out/graph.json (existing graph has more nodes; #479).')
    raise SystemExit(1)
print('Report updated with community labels successfully')
"

# BREAK


$(cat graphify-out/.graphify_python) -m graphify export html

# BREAK


$(cat graphify-out/.graphify_python) -m graphify benchmark

# BREAK

$(cat graphify-out/.graphify_python) -c "
import json
from pathlib import Path
from datetime import datetime, timezone
from graphify.detect import save_manifest

detect = json.loads(Path('graphify-out/.graphify_detect.json').read_text(encoding=\"utf-8\"))
extract = json.loads(Path('graphify-out/.graphify_extract.json').read_text(encoding=\"utf-8\"))

from graphify.cli import _stamped_manifest_files
_corpus = detect.get('all_files') or detect['files']
_manifest_files = _stamped_manifest_files(_corpus, extract, Path('.'))
_sem_types = ('document', 'paper', 'image')
_dispatched = {f for t, fl in detect['files'].items() if t in _sem_types for f in fl}
_stamped = {f for fl in _manifest_files.values() for f in fl}
_cleared = _dispatched - _stamped
_scan = {f for fl in _corpus.values() for f in fl}
save_manifest(_manifest_files, root='.', scan_corpus=_scan, clear_semantic=_cleared or None)

input_tok = extract.get('input_tokens', 0)
output_tok = extract.get('output_tokens', 0)

cost_path = Path('graphify-out/cost.json')
if cost_path.exists():
    cost = json.loads(cost_path.read_text(encoding=\"utf-8\"))
else:
    cost = {'runs': [], 'total_input_tokens': 0, 'total_output_tokens': 0}

cost['runs'].append({
    'date': datetime.now(timezone.utc).isoformat(),
    'input_tokens': input_tok,
    'output_tokens': output_tok,
    'files': detect.get('total_files', 0),
})
cost['total_input_tokens'] += input_tok
cost['total_output_tokens'] += output_tok
cost_path.write_text(json.dumps(cost, indent=2, ensure_ascii=False), encoding=\"utf-8\")

print(f'This run: {input_tok:,} input tokens, {output_tok:,} output tokens')
print(f'All time: {cost[\"total_input_tokens\"]:,} input, {cost[\"total_output_tokens\"]:,} output ({len(cost[\"runs\"])} runs)')
"
rm -f graphify-out/.graphify_detect.json graphify-out/.graphify_extract.json graphify-out/.graphify_ast.json graphify-out/.graphify_semantic.json graphify-out/.graphify_analysis.json
find graphify-out -maxdepth 1 -name '.graphify_chunk_*.json' -delete 2>/dev/null
rm -f graphify-out/.needs_update 2>/dev/null || true
