// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call, Create as $Create} from "@wailsio/runtime";

export function Greet(name: string): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2636946092, name) as any;
    return $resultPromise;
}

export function SyncCheck(timestamp: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(987921380, timestamp) as any;
    return $resultPromise;
}
