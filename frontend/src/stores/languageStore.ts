// global store for language
export class GolbalLanguage {
  static language: string = "en";

  public static getLanguage() {
    return GolbalLanguage.language;
  }

  public static setLanguage(language: string) {
    GolbalLanguage.language = language;
  }
}
