export class OrganizationConflictError extends Error {
    constructor(message: string) {
      super(message)
      this.name = 'OrganizationConflictError'
    }
  }
  