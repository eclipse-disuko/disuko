type Nullable<T> = T | null;
export class SpdxIdentifier {
  public spdxFileId = '';
  public label = '';
  public uploaded = '';
  public versionKey = '';
  public versionName = '';
  public tag = '';
  public header: Nullable<string> = null;

  constructor(spdxFileId: string, label: string, uploaded: string, versionKey: string, header: string, tag: string) {
    this.spdxFileId = spdxFileId;
    this.label = label;
    this.uploaded = uploaded;
    this.versionKey = versionKey;
    if (header.length > 0) {
      this.header = header;
    }
    this.tag = tag;
  }
}
