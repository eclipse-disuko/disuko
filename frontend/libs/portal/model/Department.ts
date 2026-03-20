export class DepartmentDto {
  public deptId = '';
  public parentDeptId = '';
  public validFrom = '';
  public descriptionEnglish = '';
  public orgAbbreviation = '';
  public skz = '';
  public companyCode = '';
  public companyName = '';
  public level = 0;
}

export class Department extends DepartmentDto {
  public constructor(dto: Department | null | undefined = null) {
    super();
    if (dto) {
      Object.assign(this, dto);
    }
  }

  public fill(dto: Department | null) {
    if (dto !== null) {
      Object.assign(this, dto);
    }
  }
}
