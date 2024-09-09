import Dataloader from 'dataloader'

const main = async () => {
    const loader = new Dataloader(
        async (keys: readonly { key: string, value: string }[]) => {
            console.warn('keys', keys)
            return keys.map(({ key }) => ({ key }))
        }
    )
    const r = await loader.load({ key: 'test', value: 'v' })
    console.warn(r)
}

main();
