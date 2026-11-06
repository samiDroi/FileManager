import { describe, expect, it } from 'vitest';
import { can } from './can';

describe('can', () => {
  it('allows a permission assigned to any of the user roles', () => {
    expect(can({
      roles: [
        { name: 'editor', permissions: ['projects.read', 'projects.create'] },
      ],
    }, 'create', 'projects')).toBe(true);
  });

  it('supports direct permissions and resource wildcards', () => {
    expect(can({ permissions: ['clients.*'] }, 'create', 'clients')).toBe(true);
    expect(can({ permissions: ['projects.read'] }, 'create', 'projects')).toBe(false);
  });

  it('allows administrators and rejects missing users or incomplete requests', () => {
    expect(can({ roles: [{ name: 'ADMIN', permissions: [] }] }, 'delete', 'users')).toBe(true);
    expect(can(null, 'read', 'files')).toBe(false);
    expect(can({ permissions: ['files.read'] }, '', 'files')).toBe(false);
  });
});
