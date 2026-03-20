export default class ErrorDialogConfig {
  public title: string;
  public titleKeyOrCode: string;
  public description: string;
  public errorCode: string;
  public reqId: string;
  public copyDesc: boolean;
  public stackTrace: any;

  public constructor() {
    this.title = '';
    this.titleKeyOrCode = '';
    this.description = '';
    this.errorCode = '';
    this.stackTrace = '';
    this.reqId = '';
    this.copyDesc = false;
  }
}
