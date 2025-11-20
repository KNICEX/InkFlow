export class PromiseGroup<T> {
  private promises: Promise<T>[] = []
  private results: T[] = []
  private errors: Error[] = []

  add(promise: Promise<T>) {
    this.promises.push(promise)
    promise
      .then((result) => {
        this.results.push(result)
      })
      .catch((error) => {
        this.errors.push(error)
      })
  }

  async wait() {
    await Promise.all(this.promises)
    return { results: this.results, errors: this.errors }
  }
}
